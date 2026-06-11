import { defineStore } from 'pinia'
import { get, post, type IApiResponse } from '@/plugins/axios'
import swal from '@/plugins/swal'
import { ref } from 'vue'

export interface IQRCode {
  id: number
  office_id: number
  office?: {
    id: number
    name: string
    address: string
    latitude: number
    longitude: number
    radius_meters: number
    is_active: boolean
  }
  code_value: string
  expires_at: string
  is_active: boolean
  created_by: number
  revoked_at?: string
  created_at: string
}

export interface IQRCodeGeneratePayload {
  office: number
  end_date: string
  end_time: string
}

export interface IQRCodeGenerateResult {
  id: number
  office_id: number
  code_value: string
  expires_at: string
  is_active: boolean
  created_at: string
}

export const useQrcodeStore = defineStore('qrcode', () => {
  const items = ref<IQRCode[]>([])
  const loading = ref<Record<string, boolean>>({
    Index: false,
    Generate: false,
    Revoke: false,
  })

  async function fetchAll(): Promise<IQRCode[]> {
    loading.value.Index = true
    try {
      const { data } = await get<IApiResponse<IQRCode[]>>('/qr-codes')
      items.value = data.data || []
      return items.value
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal memuat QR codes.'
      swal.error('Gagal', message)
      return []
    } finally {
      loading.value.Index = false
    }
  }

  async function generate(payload: IQRCodeGeneratePayload): Promise<IQRCodeGenerateResult | null> {
    loading.value.Generate = true
    try {
      const { data } = await post<IApiResponse<IQRCodeGenerateResult>>(
        '/qr-codes/generate',
        payload,
      )
      swal.success('Berhasil', data.message || 'QR code berhasil dibuat.')
      await fetchAll()
      return data.data || null
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal membuat QR code.'
      swal.error('Gagal', message)
      throw error
    } finally {
      loading.value.Generate = false
    }
  }

  async function revoke(id: number): Promise<void> {
    const result = await swal.warning(
      'Cabut QR Code',
      'Apakah Anda yakin ingin mencabut QR code ini? Karyawan tidak bisa lagi check-in menggunakan QR code ini.',
    )
    if (!result.isConfirmed) return

    loading.value.Revoke = true
    try {
      const { data } = await post<IApiResponse<any>>(`/qr-codes/${id}/revoke`, {})
      swal.success('Berhasil', data.message || 'QR code berhasil dicabut.')
      await fetchAll()
    } catch (error: any) {
      const message = error?.response?.data?.message || 'Gagal mencabut QR code.'
      swal.error('Gagal', message)
      throw error
    } finally {
      loading.value.Revoke = false
    }
  }

  return {
    items,
    loading,
    fetchAll,
    generate,
    revoke,
  }
})
