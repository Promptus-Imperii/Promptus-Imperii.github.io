const fs = require('fs').promises;
const path = require('path');
const sharp = require('sharp');

// Define the maximum dimensions
const MAX_WIDTH = 1920;
const MAX_HEIGHT = 1080;

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
        const metadata = await sharp(filePath).metadata();

        // Create a sharp instance with EXIF data
        const sharpInstance = sharp(filePath).withMetadata();

        // Create a temporary file path for the processed image
        const tempFilePath = path.join(
            path.dirname(filePath),
            `${path.basename(filePath, path.extname(filePath))}_temp${path.extname(filePath)}`
        );

        // Check if the image needs rotation based on EXIF data
        let imageProcessed = false;  // Flag to check if the image needs to be written

        if (metadata.orientation && metadata.orientation !== 1) {
            // If rotation is needed (orientation is not "1", which means no rotation), apply rotation
            console.log(`Rotating image: ${path.basename(filePath)}`);
            await sharpInstance
                .rotate() // Rotate the image based on EXIF orientation
                .toFile(tempFilePath); // Save to a temporary file
            imageProcessed = true; // Set flag to indicate that the image was processed
            console.log(`Processed and rotated ${path.basename(filePath)} and saved to temporary file`);
        }

        // If the image is larger than the maximum dimensions, resize and crop it
        if ((metadata.width > MAX_WIDTH || metadata.height > MAX_HEIGHT) && !imageProcessed) {
            // Resize and crop the image
            console.log(`Resizing image: ${path.basename(filePath)}`);
            await sharpInstance
                .resize(MAX_WIDTH, MAX_HEIGHT, { fit: 'outside' }) // Resize and crop
                .toFile(tempFilePath); // Save to a temporary file
            imageProcessed = true; // Set flag to indicate that the image was processed
            console.log(`Processed and resized ${path.basename(filePath)} and saved to temporary file`);
        }

        // If no processing was done (no rotation, no resizing), skip writing to the file
        if (!imageProcessed) {
            console.log(`No processing needed for image: ${path.basename(filePath)}`);
            return;  // Exit the function if no changes were made
        }

        // Replace the original file with the processed file
        await fs.unlink(filePath); // Delete the original file
        await fs.rename(tempFilePath, filePath); // Rename processed file to original file name
        console.log(`Replaced original file with processed version: ${path.basename(filePath)}`);

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
