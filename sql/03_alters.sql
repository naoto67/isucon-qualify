ALTER TABLE items ADD INDEX created_at_idx(`created_at`);
ALTER TABLE items ADD INDEX seller_id_and_created_at_idx(`seller_id`, `created_at`);
ALTER TABLE items ADD INDEX buyer_id_idx(`buyer_id`);
