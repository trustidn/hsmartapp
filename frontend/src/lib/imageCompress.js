/**
 * Kompresi gambar sebelum upload - agar ukuran ringan, PWA tetap cepat, storage hemat.
 * Menggunakan Canvas API (built-in, tanpa dependensi eksternal).
 *
 * Strategi:
 * - Resize ke dimensi maksimum (cukup untuk tampilan card produk / daftar)
 * - Kompres JPEG quality 0.82 (keseimbangan kualitas vs ukuran)
 * - Format WebP jika browser support (lebih kecil)
 *
 * Target: gambar 2–5MB → ~30–80KB untuk display di POS & daftar produk
 */

/** Dimensi maksimum untuk gambar produk (cukup untuk card & daftar, 2x retina) */
const MAX_WIDTH = 400
const MAX_HEIGHT = 400
const JPEG_QUALITY = 0.82
const WEBP_QUALITY = 0.8

/**
 * Deteksi dukungan WebP (format lebih ringan dari JPEG)
 */
function supportsWebP() {
  try {
    const canvas = document.createElement('canvas')
    canvas.width = 1
    canvas.height = 1
    return canvas.toDataURL('image/webp', 0.9).startsWith('data:image/webp')
  } catch {
    return false
  }
}

/**
 * Resize & kompres gambar agar ringan sebelum upload.
 *
 * @param {File} file - File gambar dari input
 * @param {Object} options - { maxWidth, maxHeight, quality, format: 'jpeg'|'webp'|'auto' }
 * @returns {Promise<Blob>} Blob gambar yang sudah dikompres
 */
export async function compressImage(file, options = {}) {
  const maxW = options.maxWidth ?? MAX_WIDTH
  const maxH = options.maxHeight ?? MAX_HEIGHT
  const quality = options.quality ?? JPEG_QUALITY
  const format = options.format ?? (supportsWebP() ? 'webp' : 'jpeg')

  return new Promise((resolve, reject) => {
    const img = new Image()
    const url = URL.createObjectURL(file)

    img.onload = () => {
      URL.revokeObjectURL(url)
      let { width, height } = img

      if (width > maxW || height > maxH) {
        const ratio = Math.min(maxW / width, maxH / height)
        width = Math.round(width * ratio)
        height = Math.round(height * ratio)
      }

      const canvas = document.createElement('canvas')
      canvas.width = width
      canvas.height = height
      const ctx = canvas.getContext('2d')
      ctx.drawImage(img, 0, 0, width, height)

      const mime = format === 'webp' ? 'image/webp' : 'image/jpeg'
      const q = format === 'webp' ? Math.min(quality, 0.9) : quality

      canvas.toBlob(
        (blob) => {
          if (blob) resolve(blob)
          else reject(new Error('Gagal kompresi gambar'))
        },
        mime,
        q
      )
    }

    img.onerror = () => {
      URL.revokeObjectURL(url)
      reject(new Error('Gagal memuat gambar'))
    }

    img.src = url
  })
}

/**
 * Kompres dan kembalikan File (untuk FormData upload).
 * Nama file dipertahankan dengan ekstensi yang sesuai format output.
 */
export async function compressImageToFile(file, options = {}) {
  const blob = await compressImage(file, options)
  const format = options.format ?? (supportsWebP() ? 'webp' : 'jpeg')
  const ext = format === 'webp' ? '.webp' : '.jpg'
  const baseName = file.name.replace(/\.[^.]+$/, '')
  return new File([blob], `${baseName}${ext}`, { type: blob.type })
}
