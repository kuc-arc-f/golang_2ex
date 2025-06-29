CREATE TABLE todos (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL,
  content TEXT,
  content_type TEXT,
  is_public INTEGER,
  food_orange INTEGER,
  food_apple INTEGER,
  food_banana INTEGER,
  food_melon INTEGER,
  food_grape INTEGER,
  created_at TEXT,
  updated_at TEXT
);