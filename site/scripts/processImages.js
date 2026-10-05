const fs = require('fs').promises;
const path = require('path');
const sharp = require('sharp');

// Define the maximum dimensions
const MAX_WIDTH_LANDSCAPE = 1920;
const MAX_HEIGHT_LANDSCAPE = 1080;

// Invert Width and Height for portrait
const MAX_WIDTH_PORTRAIT = MAX_HEIGHT_LANDSCAPE;
const MAX_HEIGHT_PORTRAIT = MAX_WIDTH_LANDSCAPE;

// Path to the assets folder
const assetsFolder = path.join(__dirname, '..', 'assets/images');

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
    }
}

// Function to process a single image
async function processImage(filePath) {
    try {
        // Create a sharp instance with EXIF data
        let sharpInstance = sharp(filePath).withMetadata();

        // Create a temporary file path for the processed image
        const tempFilePath = path.join(
            path.dirname(filePath),
            `${path.basename(filePath, path.extname(filePath))}_temp${path.extname(filePath)}`
        );

        const replaceFile = async (action) => {
            // Replace the original file with the processed file
            await fs.unlink(filePath); // Delete the original file
            await fs.rename(tempFilePath, filePath); // Rename processed file to original file name
            console.log(`Replaced original file with ${action} version: ${path.basename(filePath)}`);
        }

        // Check if the image needs rotation based on EXIF data
        if (sharpInstance.metadata.orientation && sharpInstance.metadata.orientation !== 1) {
            // If rotation is needed (orientation is not "1", which means no rotation), apply rotation
            console.log(`Rotating image: ${path.basename(filePath)}`);
            await sharpInstance
                .rotate() // Rotate the image based on EXIF orientation
                .toFile(tempFilePath); // Save to a temporary file
            await replaceFile("rotated")
            //Refresh sharpInstance to load new metadata
            sharpInstance = sharp(filePath).withMetadata();
            console.log(`Processed and rotated ${path.basename(filePath)}`);
        }

        const resizeImage = async (width, height) => {
            // Resize and crop the image
            console.log(`Resizing image: ${path.basename(filePath)}`);
            await sharpInstance
                .resize(width, height, { fit: 'outside' }) // Resize and crop
                .toFile(tempFilePath); // Save to a temporary file
            await replaceFile("resized")
            console.log(`Processed and resized ${path.basename(filePath)}`);
        }

        // If landscape mode and larger than 1080p, resize it
        if (
            sharpInstance.metadata.width > sharpInstance.metadata.height &&
            sharpInstance.metadata.width > MAX_WIDTH_LANDSCAPE || sharpInstance.metadata.height > MAX_HEIGHT_LANDSCAPE
        ) {
            // Image is landscape mode
            resizeImage(MAX_WIDTH_LANDSCAPE, MAX_HEIGHT_LANDSCAPE)
        }
        // If portrait mode and larger than 1080p
        if(
            sharpInstance.metadata.width < sharpInstance.metadata.height &&
            sharpInstance.metadata.width > MAX_WIDTH_PORTRAIT || sharpInstance.metadata.height > MAX_HEIGHT_PORTRAIT
        ) {
            resizeImage(MAX_WIDTH_PORTRAIT, MAX_HEIGHT_PORTRAIT)
        }

    } catch (error) {
        console.error(`Error processing ${path.basename(filePath)}:`, error);
    }
}



// Helper function to check if a file is an image
// HEIC cannot be processed
function isImage(filePath) {
    const imageExtensions = ['.jpg', '.jpeg', '.png', '.webp'];
    const ext = path.extname(filePath).toLowerCase();
    return imageExtensions.includes(ext);
}

// Start the recursive processing
(async () => {
    if (!(await fs.stat(assetsFolder).catch(() => false))) {
        console.error(`Folder "${assetsFolder}" does not exist.`);
        process.exit(1);
    }
    await processFolder(assetsFolder);
})();
