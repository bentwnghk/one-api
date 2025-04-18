export const priceType = [
  { value: 'tokens', label: '按Token收費' },
  { value: 'times', label: '按次收費' }
];

export function ValueFormatter(value) {
  if (value == null) {
    return '';
  }
  if (value === 0) {
    return 'Free';
  }
  return `$${parseFloat(value * 0.002).toFixed(6)} / ￥${parseFloat(value * 0.014).toFixed(6)}`;
}
