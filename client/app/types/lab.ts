export interface LabResult {
  id: string
  client_id: string
  nutritionist_id: string
  title: string
  result_type: string
  test_date: string | null
  original_name: string | null
  file_type: string | null
  file_size: number
  link: string | null
  notes: string
  created_at: string
}

export interface UploadLabResultRequest {
  title: string
  result_type: string
  test_date?: string
  notes?: string
  file?: File
  link?: string
}

export type LabResourceType = 'file' | 'link'

export function getLabResourceType(result: LabResult): LabResourceType {
  return result.link !== null ? 'link' : 'file'
}
