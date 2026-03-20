INSERT INTO sports (name, slug) VALUES
  ('Football', 'football'),
  ('Ice Hockey', 'ice-hockey'),
  ('Basketball', 'basketball'),
  ('Volleyball', 'volleyball'),
  ('Handball', 'handball'),
  ('Tennis', 'tennis')
ON CONFLICT (slug) DO NOTHING;