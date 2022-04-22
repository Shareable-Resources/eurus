-- Deploy eurus-backend-db-user:create_user_kyc_images to pg
BEGIN;

SET search_path to public;

-- XXX Add DDLs here.
CREATE TABLE IF NOT EXISTS user_kyc_images (
    user_kyc_id BIGINT NOT NULL,
    doc_type SMALLINT NOT NULL,
    image_seq INT NOT NULL,
    status SMALLINT DEFAULT 0 NOT NULL,
    image_path TEXT,
    created_date TIMESTAMP WITH TIME ZONE NOT NULL,
    last_modified_date TIMESTAMP WITH TIME ZONE NOT NULL,
    reject_reason VARCHAR(255),
    operator_id BIGINT,
    PRIMARY KEY (user_kyc_id,doc_type,image_seq)
);

COMMENT ON COLUMN user_kyc_images.user_kyc_id IS 'Foreign key to represent a kyc status';
COMMENT ON COLUMN user_kyc_images.doc_type IS 'Doc Type can be 0(IdFront),1(IdBack),2(Passport),3(Selfie)';
COMMENT ON COLUMN user_kyc_images.image_seq IS 'Image Seq represents how many times the images has been uploaded';
COMMENT ON COLUMN user_kyc_images.status IS 'Status can be 0(Received), 1(Uploaded to S3 Server), 2(Approved), 3(Rejected)';
COMMENT ON COLUMN user_kyc_images.image_path IS 'The url link for the image. This link will be updated once KYC Server upload the image to S3 Server';
COMMENT ON COLUMN user_kyc_images.created_date IS 'Created date';
COMMENT ON COLUMN user_kyc_images.last_modified_date IS 'Last Modified Date';
COMMENT ON COLUMN user_kyc_images.reject_reason IS 'Reject Reason of the uploaded image, entered by CS admin';
COMMENT ON COLUMN user_kyc_images.operator_id IS 'The CS admin who approve/reject or do things to this record';

COMMIT;
