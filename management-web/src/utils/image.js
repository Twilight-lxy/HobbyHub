/**
 * 图片工具函数
 */

// 图片服务器基础URL
const BASE_IMAGE_URL = 'http://localhost:8081'

/**
 * 格式化图片URL，添加基础URL前缀
 * @param {string} url 图片路径
 * @param {string} defaultUrl 默认图片路径
 * @returns {string} 完整的图片URL
 */
export function formatImageUrl(url, defaultUrl = '') {
  if (!url) return defaultUrl
  
  // 如果已经是完整URL，则直接返回
  if (url.startsWith('http://') || url.startsWith('https://')) {
    return url
  }
  
  // 添加基础URL前缀
  return `${BASE_IMAGE_URL}${url}`
}

/**
 * 格式化图片数组中的第一张图片
 * @param {Array} images 图片数组
 * @param {string} defaultUrl 默认图片路径
 * @returns {string} 第一张图片的完整URL
 */
export function formatFirstImage(images, defaultUrl = '') {
  if (!images || !images.length) return defaultUrl
  return formatImageUrl(images[0], defaultUrl)
}

export default {
  formatImageUrl,
  formatFirstImage
} 