-- Migration: Add rejected_fields column to kyc_submissions
-- This column stores an array of field names that were flagged as incorrect during rejection.
-- Run this on the PostgreSQL database before deploying the updated backend.

ALTER TABLE kyc_submissions
ADD COLUMN IF NOT EXISTS rejected_fields TEXT[] DEFAULT NULL;

COMMENT ON COLUMN kyc_submissions.rejected_fields IS 'Array of field names flagged as incorrect during rejection (e.g., full_name, nik, address, birthdate, phone, ktp_image)';
