export const RANK_SCORE_COLOR_BANDS = Object.freeze([
  { min: 0, max: 699, label: '0-699', color: '#B4B4B4' },
  { min: 700, max: 1099, label: '700-1099', color: '#C177E7' },
  { min: 1100, max: 1499, label: '1100-1499', color: '#5EA1F4' },
  { min: 1500, max: 1999, label: '1500-1999', color: '#25BB9B' },
  { min: 2000, max: 2399, label: '2000-2399', color: '#E8CA2D' },
  { min: 2400, max: 2799, label: '2400-2799', color: '#FC7123' },
  { min: 2800, max: Infinity, label: '2800+', color: '#FF0000' }
]);

export function getRankScoreColor(score) {
  const number = Number(score);
  const safeScore = Number.isFinite(number) ? number : 0;
  const band = RANK_SCORE_COLOR_BANDS.find(item => safeScore >= item.min && safeScore <= item.max);

  return band?.color || RANK_SCORE_COLOR_BANDS[0].color;
}
