const fs = require('fs');
const path = require('path');

const baseDir = path.join(__dirname, '../public/downloaded_images');
outputPath = path.join(__dirname, '../src/_dataImage/');

const library = {};


// Loop semua folder tier_*
fs.readdirSync(baseDir).forEach(tierDir => {
  if (tierDir.startsWith('tier_')) {
    const tier = parseInt(tierDir.split('_')[1]);
    const tierPath = path.join(baseDir, tierDir);

    fs.readdirSync(tierPath).forEach(file => {
      filename = file.split('_');
      elementName = '';
      for (let i = 0; i < filename.length - 1; i++) { 
        if(i == filename.length - 2){
            elementName += filename[i];
            break;
        }
        elementName += filename[i] + ' ';
      }
      library[elementName] = {
        svgPath: `/downloaded_images/${tierDir}/${file}`,
        tier: tier
      };
    });
  }
});

if (!fs.existsSync(outputPath)) {
  fs.mkdirSync(outputPath, { recursive: true });
}

const sortedLibrary = Object.keys(library)
  .sort((a, b) => a.localeCompare(b))
  .reduce((acc, key) => {
    acc[key] = library[key];
    return acc;
  }, {});


outputPath += "elementLibrary.json"
fs.writeFileSync(outputPath, JSON.stringify(sortedLibrary, null, 2));
console.log("Library generated at:", outputPath);