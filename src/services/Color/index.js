/**
 * Available pastel colors.
 * @type {Array}
 */
const availableColors = [
  '#006E6D',
  '#2A4B7C',
  '#3F69AA',
  '#77212E',
  '#577284',
  '#6C4F3D',
  '#797B3A',
  '#935529',
  '#BD3D3A',
  '#9B1B30',
  '#E08119',
  '#6B5B95',
  '#F96714',
  '#485167',
  '#2E4A62',
  '#264E36',
];

let colorIndex;
let colorMemory;

export function clear() {
  colorMemory = {};
  colorIndex = 0;
}

export function get(str) {
  let color = colorMemory[str];
  if (!color) {
    color = colorMemory[str] =
      availableColors[colorIndex++ % availableColors.length];
  }

  return color;
}
