import jsPDF from 'jspdf'

function formatNum(n) {
  return Number(n || 0).toLocaleString('id-ID')
}

function formatDate(s) {
  if (!s) return ''
  return new Date(s).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' })
}

function formatDateTime(s) {
  if (!s) return ''
  return new Date(s).toLocaleString('id-ID', { day: '2-digit', month: 'short', hour: '2-digit', minute: '2-digit' })
}

/** Ekspor laporan ke PDF */
export function exportReportPdf(summary, transactions, expenses, filterLabel, settings) {
  const doc = new jsPDF()
  const name = settings?.name || 'Laporan'
  let y = 15

  doc.setFontSize(16)
  doc.text(name, 14, y)
  y += 6
  doc.setFontSize(10)
  doc.text(`Laporan ${filterLabel}`, 14, y)
  doc.text(formatDate(new Date()), 190, y, { align: 'right' })
  y += 12

  doc.setDrawColor(220, 220, 220)
  doc.line(14, y, 196, y)
  y += 10

  doc.setFontSize(11)
  doc.text('Ringkasan', 14, y)
  y += 8
  doc.setFontSize(10)
  doc.text(`Penjualan: Rp ${formatNum(summary?.sales_total ?? 0)}`, 14, y)
  y += 6
  doc.text(`Pengeluaran: Rp ${formatNum(summary?.expense_total ?? 0)}`, 14, y)
  y += 6
  doc.text(`Untung: Rp ${formatNum(summary?.profit ?? 0)}`, 14, y)
  doc.text(`Transaksi: ${summary?.transactions ?? 0}`, 100, y)
  y += 12

  if (transactions?.length) {
    if (y > 250) { doc.addPage(); y = 15 }
    doc.setFontSize(11)
    doc.text('Transaksi', 14, y)
    y += 8
    doc.setFontSize(9)
    for (const s of transactions.slice(0, 30)) {
      if (y > 270) { doc.addPage(); y = 15 }
      doc.text(formatDateTime(s.created_at), 14, y)
      doc.text(`Rp ${formatNum(s.total)}`, 190, y, { align: 'right' })
      y += 6
    }
    if (transactions.length > 30) doc.text(`... dan ${transactions.length - 30} transaksi lainnya`, 14, y)
    y += 10
  }

  if (expenses?.length) {
    if (y > 250) { doc.addPage(); y = 15 }
    doc.setFontSize(11)
    doc.text('Pengeluaran', 14, y)
    y += 8
    doc.setFontSize(9)
    for (const e of expenses.slice(0, 30)) {
      if (y > 270) { doc.addPage(); y = 15 }
      doc.text((e.name || '-').substring(0, 35), 14, y)
      doc.text(`Rp ${formatNum(e.amount)}`, 190, y, { align: 'right' })
      y += 6
    }
  }

  doc.save(`laporan-${filterLabel.replace(/\s/g, '-')}-${new Date().toISOString().slice(0, 10)}.pdf`)
}

/** Ekspor laporan ke CSV (bisa dibuka di Excel) */
export function exportReportCsv(summary, transactions, expenses, filterLabel) {
  const rows = []
  rows.push(['Laporan', filterLabel])
  rows.push(['Tanggal', formatDate(new Date())])
  rows.push([])
  rows.push(['RINGKASAN'])
  rows.push(['Penjualan', `Rp ${formatNum(summary?.sales_total ?? 0)}`])
  rows.push(['Pengeluaran', `Rp ${formatNum(summary?.expense_total ?? 0)}`])
  rows.push(['Untung', `Rp ${formatNum(summary?.profit ?? 0)}`])
  rows.push(['Transaksi', summary?.transactions ?? 0])
  rows.push([])
  rows.push(['TRANSAKSI', 'Tanggal', 'Total', 'Metode'])
  for (const s of transactions || []) {
    rows.push(['', formatDateTime(s.created_at), s.total ?? 0, s.payment_method || 'cash'])
  }
  rows.push([])
  rows.push(['PENGELUARAN', 'Nama', 'Jumlah', 'Tanggal'])
  for (const e of expenses || []) {
    rows.push(['', e.name || '', e.amount ?? 0, formatDateTime(e.created_at)])
  }

  const csv = rows.map((r) => r.map((c) => `"${String(c).replace(/"/g, '""')}"`).join(',')).join('\n')
  const bom = '\uFEFF'
  const blob = new Blob([bom + csv], { type: 'text/csv;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `laporan-${filterLabel.replace(/\s/g, '-')}-${new Date().toISOString().slice(0, 10)}.csv`
  a.click()
  URL.revokeObjectURL(url)
}
