// src/utils/date.js

export function formatDate(dateStr) {
  console.log("formatDate")
  return new Date(dateStr).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  });
}