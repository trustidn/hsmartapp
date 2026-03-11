import jsPDF from 'jspdf'

function formatNum(n) {
  return Number(n || 0).toLocaleString('id-ID')
}

function formatDate(s) {
  if (!s) return ''
  return new Date(s).toLocaleString('id-ID')
}

function paymentLabel(m) {
  const map = { cash: 'Tunai', qris: 'QRIS', transfer: 'Transfer', ewallet: 'E-Wallet' }
  return map[m] || m || 'Tunai'
}

export function getReceiptWhatsAppText(sale, settings) {
  const name = settings?.name || 'HSmart POS'
  let text = `Terima kasih telah berbelanja.\n\n${name}\n${formatDate(sale.created_at)}\n\n`
  for (const i of sale.items || []) {
    const name_ = i.product_name || i.name || 'Item'
    const qty = i.qty || 1
    const subtotal = i.subtotal ?? (i.price || 0) * qty
    text += `${name_} x${qty} = Rp ${formatNum(subtotal)}\n`
  }
  text += `\nTotal: Rp ${formatNum(sale.total)}\n`
  if (sale.amount_paid != null && sale.amount_paid > 0) {
    text += `Jumlah bayar: Rp ${formatNum(sale.amount_paid)}\n`
    if (sale.change != null && sale.change > 0) {
      text += `Kembalian: Rp ${formatNum(sale.change)}\n`
    }
  }
  text += `Metode: ${paymentLabel(sale.payment_method)}`
  if (settings?.receipt_footer) text += `\n\n${settings.receipt_footer}`
  return text
}

export function exportReceiptPdf(sale, settings) {
  const doc = new jsPDF({ format: [80, 200], unit: 'mm' })
  const name = settings?.name || 'HSmart POS'
  let y = 10

  doc.setFontSize(12)
  doc.text(name, 40, y, { align: 'center' })
  y += 6
  doc.setFontSize(8)
  doc.text(formatDate(sale.created_at), 40, y, { align: 'center' })
  y += 8

  doc.setDrawColor(200, 200, 200)
  doc.line(5, y, 75, y)
  y += 5

  for (const i of sale.items || []) {
    const name_ = (i.product_name || i.name || 'Item').substring(0, 25)
    const qty = i.qty || 1
    const subtotal = i.subtotal ?? (i.price || 0) * qty
    doc.text(`${name_} x${qty}`, 5, y)
    doc.text(`Rp ${formatNum(subtotal)}`, 75, y, { align: 'right' })
    y += 5
  }

  y += 3
  doc.line(5, y, 75, y)
  y += 5
  doc.setFont(undefined, 'bold')
  doc.text('Total', 5, y)
  doc.text(`Rp ${formatNum(sale.total)}`, 75, y, { align: 'right' })
  y += 5
  doc.setFont(undefined, 'normal')
  if (sale.amount_paid != null && sale.amount_paid > 0) {
    doc.text('Jumlah bayar', 5, y)
    doc.text(`Rp ${formatNum(sale.amount_paid)}`, 75, y, { align: 'right' })
    y += 5
    if (sale.change != null && sale.change > 0) {
      doc.setFont(undefined, 'bold')
      doc.text('Kembalian', 5, y)
      doc.text(`Rp ${formatNum(sale.change)}`, 75, y, { align: 'right' })
      y += 5
      doc.setFont(undefined, 'normal')
    }
  }
  doc.text(`Metode: ${paymentLabel(sale.payment_method)}`, 5, y)
  y += 8

  if (settings?.receipt_footer) {
    doc.setFontSize(6)
    doc.text(settings.receipt_footer, 40, y, { align: 'center' })
  }

  doc.save(`struk-${sale.id || Date.now()}.pdf`)
}
