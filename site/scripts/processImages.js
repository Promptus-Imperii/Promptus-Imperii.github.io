const fs = require('fs').promises;
const path = require('path');
const sharp = require('sharp');

// Define the maximum dimensions
const MAX_WIDTH_LANDSCAPE = 1920;
const MAX_HEIGHT_LANDSCAPE = 1080;

// Invert Width and Height for portrait
const MAX_WIDTH_PORTRAIT = MAX_HEIGHT_LANDSCAPE;
const MAX_HEIGHT_PORTRAIT = MAX_WIDTH_LANDSCAPE;

// Resizing or rotating always encodes a JPEG again. 90 loses nothing visible; Hugo makes its own
// smaller copies for the website anyway.
const JPEG_QUALITY = 90;

// Path to the assets folder
const assetsFolder = path.join(__dirname, '..', 'assets/images');

// Set when any folder or image fails, so the script exits with an error code.
let failed = false;

// Recursive function to process files in directories
async function processFolder(folder) {
    try {
        const entries = await fs.readdir(folder, { withFileTypes: true });
        for (const entry of entries) {
            const fullPath = path.join(folder, entry.name);
            if (entry.isDirectory()) {
                // If entry is a directory, recurse into it
                await processFolder(fullPath);
            } else if (isImage(fullPath)) {
                // If entry is an image, process it
                await processImage(fullPath);
            } else {
                console.log(`${entry.name} is not an image, skipping.`);
            }
        }
    } catch (err) {
        console.error(`Error reading folder ${folder}:`, err);
        failed = true;
    }
}

// Rotates an image according to its EXIF orientation and scales it down to fit within 1080p.
async function processImage(filePath) {
    try {
        const { width, height, orientation = 1, chromaSubsampling } = await sharp(filePath).metadata();
        // Orientations 5 to 8 turn the image by 90°, so the displayed width and height are swapped.
        const [displayWidth, displayHeight] = orientation >= 5 ? [height, width] : [width, height];
        const [maxWidth, maxHeight] = displayWidth >= displayHeight
            ? [MAX_WIDTH_LANDSCAPE, MAX_HEIGHT_LANDSCAPE]
            : [MAX_WIDTH_PORTRAIT, MAX_HEIGHT_PORTRAIT];
        const rotate = orientation !== 1;
        const resize = displayWidth > maxWidth || displayHeight > maxHeight;
        if (!rotate && !resize) {
            return;
        }

        // autoOrient turns the pixels and resets the EXIF orientation, so every viewer shows the same thing.
        let image = sharp(filePath).autoOrient().keepMetadata();
        if (resize) {
            image = image.resize(maxWidth, maxHeight, { fit: 'inside' });
        }
        const ext = path.extname(filePath).toLowerCase();
        if (ext === '.jpg' || ext === '.jpeg') {
            image = image.jpeg({
                quality: JPEG_QUALITY,
                // Keep full colour detail if the original has it (sharp would reduce it to 4:2:0).
                chromaSubsampling: chromaSubsampling === '4:4:4' ? '4:4:4' : '4:2:0',
            });
        } else if (ext === '.png') {
            // sharp's default PNG compression is lighter than most tools use, which can make files grow.
            image = image.png({ compressionLevel: 9, adaptiveFiltering: true });
        }
        // Write next to the original first, so a failure never leaves a half-written image behind.
        const tempFilePath = tempPath(filePath);
        try {
            const { size } = await image.toFile(tempFilePath);
            // A resize that only makes the file bigger isn't worth it (this happens when the original
            // was saved at a lower quality). A rotation is always kept: otherwise the image shows sideways.
            if (!rotate && size > (await fs.stat(filePath)).size) {
                console.log(`Kept ${path.basename(filePath)}: the resized version was larger than the original`);
                return;
            }
            await fs.rename(tempFilePath, filePath);
        } finally {
            await fs.rm(tempFilePath, { force: true });
        }
        console.log(`Processed ${path.basename(filePath)}:${rotate ? ' rotated' : ''}${resize ? ' resized' : ''}`);
    } catch (error) {
        console.error(`Error processing ${path.basename(filePath)}:`, error);
        failed = true;
    }
}

function tempPath(filePath) {
    const ext = path.extname(filePath);
    return path.join(path.dirname(filePath), `${path.basename(filePath, ext)}_temp${ext}`);
}

// Helper function to check if a file is an image
// HEIC cannot be processed
function isImage(filePath) {
    const imageExtensions = ['.jpg', '.jpeg', '.png', '.webp'];
    const ext = path.extname(filePath).toLowerCase();
    // Skip temporary files from an interrupted earlier run.
    return imageExtensions.includes(ext) && !path.basename(filePath, path.extname(filePath)).endsWith('_temp');
}

// Start the recursive processing
(async () => {
    if (!(await fs.stat(assetsFolder).catch(() => false))) {
        console.error(`Folder "${assetsFolder}" does not exist.`);
        process.exit(1);
    }
    await processFolder(assetsFolder);
    if (failed) {
        process.exit(1);
    }
})();
