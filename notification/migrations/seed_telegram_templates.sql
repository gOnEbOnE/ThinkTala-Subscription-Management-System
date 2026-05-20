-- Seed: Template Telegram untuk broadcast news/promo/market analysis
-- Jalankan sekali di database setelah migration add_telegram_chat_id.sql

INSERT INTO notification_templates (id, name, event_type, channel, subject, content, created_by)
VALUES
  (
    gen_random_uuid(),
    'Telegram - News Broadcast',
    'news_broadcast',
    'telegram',
    NULL,
    E'📢 *{{title}}*\n\n{{content}}\n\n_Halo {{name}}, informasi ini dikirim khusus untukmu dari Propensuy._',
    'system'
  ),
  (
    gen_random_uuid(),
    'Telegram - Promo Notification',
    'news_promo',
    'telegram',
    NULL,
    E'🔥 *Promo Spesial: {{title}}*\n\n{{content}}\n\n👉 Cek sekarang di aplikasi Propensuy!',
    'system'
  ),
  (
    gen_random_uuid(),
    'Telegram - Market Analysis',
    'news_analysis',
    'telegram',
    NULL,
    E'📊 *Market Analysis: {{title}}*\n\n{{content}}\n\n_Dikirim otomatis saat analisis baru terbit._',
    'system'
  )
ON CONFLICT DO NOTHING;

-- Daftarkan event_type ke registry
INSERT INTO notification_event_types (event_type)
VALUES
  ('news_broadcast'),
  ('news_promo'),
  ('news_analysis')
ON CONFLICT (event_type) DO NOTHING;
