create table todos(
  "id" SERIAL NOT NULL,
  title TEXT,
  content TEXT,
  completed BOOLEAN NOT NULL DEFAULT false,
  content_type TEXT,
  is_public BOOLEAN NOT NULL DEFAULT false,
  food_orange BOOLEAN NOT NULL DEFAULT false,
  food_apple BOOLEAN NOT NULL DEFAULT false,
  food_banana BOOLEAN NOT NULL DEFAULT false,
  food_melon BOOLEAN NOT NULL DEFAULT false,
  food_grape BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT "Todo_pkey" PRIMARY KEY ("id")
);