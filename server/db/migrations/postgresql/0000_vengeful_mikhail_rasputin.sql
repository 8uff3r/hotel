CREATE TABLE "accounts" (
	"id" serial PRIMARY KEY NOT NULL,
	"hotel_id" integer,
	"account_code" text NOT NULL,
	"account_name" text NOT NULL,
	"account_type" text NOT NULL,
	"account_sub_type" text,
	"parent_id" integer,
	"is_active" boolean DEFAULT true NOT NULL,
	"is_system" boolean DEFAULT false NOT NULL,
	"description" text,
	"normal_balance" text DEFAULT 'debit' NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE "expenses" (
	"id" serial PRIMARY KEY NOT NULL,
	"hotel_id" integer,
	"expense_date" timestamp NOT NULL,
	"description" text NOT NULL,
	"amount" numeric(15, 2) NOT NULL,
	"category" text NOT NULL,
	"vendor" text,
	"reference" text,
	"payment_method" text,
	"payment_status" text DEFAULT 'pending' NOT NULL,
	"account_id" integer,
	"receipt_number" text,
	"notes" text,
	"created_by" integer,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE "guests" (
	"id" serial PRIMARY KEY NOT NULL,
	"first_name" text NOT NULL,
	"last_name" text NOT NULL,
	"email" text,
	"phone" text,
	"id_type" text,
	"id_number" text,
	"address" text,
	"city" text,
	"country" text,
	"notes" text,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE "hotels" (
	"id" serial PRIMARY KEY NOT NULL,
	"name" text NOT NULL,
	"address" text NOT NULL,
	"phone" text,
	"email" text,
	"logo_url" text,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE "income" (
	"id" serial PRIMARY KEY NOT NULL,
	"hotel_id" integer,
	"income_date" timestamp NOT NULL,
	"description" text NOT NULL,
	"amount" numeric(15, 2) NOT NULL,
	"category" text NOT NULL,
	"source" text,
	"reference" text,
	"payment_method" text,
	"payment_status" text DEFAULT 'pending' NOT NULL,
	"account_id" integer,
	"reservation_id" integer,
	"notes" text,
	"created_by" integer,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE "journal_entries" (
	"id" serial PRIMARY KEY NOT NULL,
	"hotel_id" integer,
	"entry_number" text NOT NULL,
	"entry_date" timestamp NOT NULL,
	"description" text NOT NULL,
	"reference" text,
	"status" text DEFAULT 'draft' NOT NULL,
	"total_debit" numeric(15, 2) DEFAULT '0' NOT NULL,
	"total_credit" numeric(15, 2) DEFAULT '0' NOT NULL,
	"created_by" integer,
	"posted_at" timestamp,
	"voided_at" timestamp,
	"void_reason" text,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE "journal_lines" (
	"id" serial PRIMARY KEY NOT NULL,
	"entry_id" integer NOT NULL,
	"account_id" integer NOT NULL,
	"description" text,
	"debit" numeric(15, 2) DEFAULT '0' NOT NULL,
	"credit" numeric(15, 2) DEFAULT '0' NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE "parking_lots" (
	"id" serial PRIMARY KEY NOT NULL,
	"hotel_id" integer,
	"name" text NOT NULL,
	"location" text,
	"total_spots" integer DEFAULT 0 NOT NULL,
	"hourly_rate" text DEFAULT '0' NOT NULL,
	"daily_rate" text DEFAULT '0' NOT NULL,
	"status" text DEFAULT 'active' NOT NULL,
	"description" text,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE "parking_spots" (
	"id" serial PRIMARY KEY NOT NULL,
	"lot_id" integer,
	"spot_number" text NOT NULL,
	"floor" text,
	"spot_type" text DEFAULT 'standard' NOT NULL,
	"status" text DEFAULT 'available' NOT NULL,
	"is_covered" boolean DEFAULT false NOT NULL,
	"description" text,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE "parking_transactions" (
	"id" serial PRIMARY KEY NOT NULL,
	"lot_id" integer,
	"spot_id" integer,
	"guest_id" integer,
	"reservation_id" integer,
	"license_plate" text NOT NULL,
	"entry_time" timestamp NOT NULL,
	"exit_time" timestamp,
	"hours_parked" real,
	"rate_applied" text,
	"amount_due" real DEFAULT 0 NOT NULL,
	"amount_paid" real DEFAULT 0 NOT NULL,
	"status" text DEFAULT 'active' NOT NULL,
	"payment_status" text DEFAULT 'pending' NOT NULL,
	"payment_method" text,
	"notes" text,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE "reservations" (
	"id" serial PRIMARY KEY NOT NULL,
	"hotel_id" integer,
	"guest_id" integer NOT NULL,
	"room_id" integer NOT NULL,
	"check_in_date" timestamp NOT NULL,
	"check_out_date" timestamp NOT NULL,
	"actual_check_in" timestamp,
	"actual_check_out" timestamp,
	"status" text DEFAULT 'pending' NOT NULL,
	"total_amount" real DEFAULT 0 NOT NULL,
	"paid_amount" real DEFAULT 0 NOT NULL,
	"payment_status" text DEFAULT 'pending' NOT NULL,
	"special_requests" text,
	"number_of_guests" integer DEFAULT 1 NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL,
	"created_by" integer
);
--> statement-breakpoint
CREATE TABLE "room_types" (
	"id" serial PRIMARY KEY NOT NULL,
	"name" text NOT NULL,
	"description" text,
	"base_price_multiplier" real DEFAULT 1 NOT NULL,
	"amenities" text,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL,
	CONSTRAINT "room_types_name_unique" UNIQUE("name")
);
--> statement-breakpoint
CREATE TABLE "rooms" (
	"id" serial PRIMARY KEY NOT NULL,
	"hotel_id" integer,
	"room_number" text NOT NULL,
	"room_type" text DEFAULT 'single' NOT NULL,
	"floor" integer,
	"capacity" integer DEFAULT 2 NOT NULL,
	"base_price" real DEFAULT 0 NOT NULL,
	"status" text DEFAULT 'available' NOT NULL,
	"amenities" text,
	"description" text,
	"images" text,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE "user_roles" (
	"id" serial PRIMARY KEY NOT NULL,
	"user_id" integer,
	"role" text NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE "users" (
	"id" serial PRIMARY KEY NOT NULL,
	"email" text NOT NULL,
	"password_hash" text NOT NULL,
	"first_name" text NOT NULL,
	"last_name" text NOT NULL,
	"role" text DEFAULT 'staff' NOT NULL,
	"avatar" text,
	"is_active" boolean DEFAULT true NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL,
	CONSTRAINT "users_email_unique" UNIQUE("email")
);
--> statement-breakpoint
CREATE TABLE "vehicles" (
	"id" serial PRIMARY KEY NOT NULL,
	"guest_id" integer,
	"license_plate" text NOT NULL,
	"vehicle_type" text DEFAULT 'car' NOT NULL,
	"make" text,
	"model" text,
	"color" text,
	"is_registered" integer DEFAULT 1 NOT NULL,
	"notes" text,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
ALTER TABLE "accounts" ADD CONSTRAINT "accounts_hotel_id_hotels_id_fk" FOREIGN KEY ("hotel_id") REFERENCES "public"."hotels"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "expenses" ADD CONSTRAINT "expenses_hotel_id_hotels_id_fk" FOREIGN KEY ("hotel_id") REFERENCES "public"."hotels"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "expenses" ADD CONSTRAINT "expenses_account_id_accounts_id_fk" FOREIGN KEY ("account_id") REFERENCES "public"."accounts"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "expenses" ADD CONSTRAINT "expenses_created_by_users_id_fk" FOREIGN KEY ("created_by") REFERENCES "public"."users"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "income" ADD CONSTRAINT "income_hotel_id_hotels_id_fk" FOREIGN KEY ("hotel_id") REFERENCES "public"."hotels"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "income" ADD CONSTRAINT "income_account_id_accounts_id_fk" FOREIGN KEY ("account_id") REFERENCES "public"."accounts"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "income" ADD CONSTRAINT "income_created_by_users_id_fk" FOREIGN KEY ("created_by") REFERENCES "public"."users"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "journal_entries" ADD CONSTRAINT "journal_entries_hotel_id_hotels_id_fk" FOREIGN KEY ("hotel_id") REFERENCES "public"."hotels"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "journal_entries" ADD CONSTRAINT "journal_entries_created_by_users_id_fk" FOREIGN KEY ("created_by") REFERENCES "public"."users"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "journal_lines" ADD CONSTRAINT "journal_lines_entry_id_journal_entries_id_fk" FOREIGN KEY ("entry_id") REFERENCES "public"."journal_entries"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "journal_lines" ADD CONSTRAINT "journal_lines_account_id_accounts_id_fk" FOREIGN KEY ("account_id") REFERENCES "public"."accounts"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "parking_lots" ADD CONSTRAINT "parking_lots_hotel_id_hotels_id_fk" FOREIGN KEY ("hotel_id") REFERENCES "public"."hotels"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "parking_spots" ADD CONSTRAINT "parking_spots_lot_id_parking_lots_id_fk" FOREIGN KEY ("lot_id") REFERENCES "public"."parking_lots"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "parking_transactions" ADD CONSTRAINT "parking_transactions_lot_id_parking_lots_id_fk" FOREIGN KEY ("lot_id") REFERENCES "public"."parking_lots"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "parking_transactions" ADD CONSTRAINT "parking_transactions_spot_id_parking_spots_id_fk" FOREIGN KEY ("spot_id") REFERENCES "public"."parking_spots"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "parking_transactions" ADD CONSTRAINT "parking_transactions_guest_id_guests_id_fk" FOREIGN KEY ("guest_id") REFERENCES "public"."guests"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "parking_transactions" ADD CONSTRAINT "parking_transactions_reservation_id_reservations_id_fk" FOREIGN KEY ("reservation_id") REFERENCES "public"."reservations"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "reservations" ADD CONSTRAINT "reservations_hotel_id_hotels_id_fk" FOREIGN KEY ("hotel_id") REFERENCES "public"."hotels"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "reservations" ADD CONSTRAINT "reservations_guest_id_guests_id_fk" FOREIGN KEY ("guest_id") REFERENCES "public"."guests"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "reservations" ADD CONSTRAINT "reservations_room_id_rooms_id_fk" FOREIGN KEY ("room_id") REFERENCES "public"."rooms"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "reservations" ADD CONSTRAINT "reservations_created_by_users_id_fk" FOREIGN KEY ("created_by") REFERENCES "public"."users"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "rooms" ADD CONSTRAINT "rooms_hotel_id_hotels_id_fk" FOREIGN KEY ("hotel_id") REFERENCES "public"."hotels"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "user_roles" ADD CONSTRAINT "user_roles_user_id_users_id_fk" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "vehicles" ADD CONSTRAINT "vehicles_guest_id_guests_id_fk" FOREIGN KEY ("guest_id") REFERENCES "public"."guests"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
CREATE INDEX "accounts_hotel_id_idx" ON "accounts" USING btree ("hotel_id");--> statement-breakpoint
CREATE INDEX "accounts_account_code_idx" ON "accounts" USING btree ("account_code");--> statement-breakpoint
CREATE INDEX "accounts_account_type_idx" ON "accounts" USING btree ("account_type");--> statement-breakpoint
CREATE INDEX "accounts_parent_id_idx" ON "accounts" USING btree ("parent_id");--> statement-breakpoint
CREATE INDEX "expenses_hotel_id_idx" ON "expenses" USING btree ("hotel_id");--> statement-breakpoint
CREATE INDEX "expenses_expense_date_idx" ON "expenses" USING btree ("expense_date");--> statement-breakpoint
CREATE INDEX "expenses_category_idx" ON "expenses" USING btree ("category");--> statement-breakpoint
CREATE INDEX "expenses_payment_status_idx" ON "expenses" USING btree ("payment_status");--> statement-breakpoint
CREATE INDEX "guests_email_idx" ON "guests" USING btree ("email");--> statement-breakpoint
CREATE INDEX "guests_phone_idx" ON "guests" USING btree ("phone");--> statement-breakpoint
CREATE INDEX "guests_name_idx" ON "guests" USING btree ("last_name","first_name");--> statement-breakpoint
CREATE INDEX "hotels_name_idx" ON "hotels" USING btree ("name");--> statement-breakpoint
CREATE INDEX "income_hotel_id_idx" ON "income" USING btree ("hotel_id");--> statement-breakpoint
CREATE INDEX "income_income_date_idx" ON "income" USING btree ("income_date");--> statement-breakpoint
CREATE INDEX "income_category_idx" ON "income" USING btree ("category");--> statement-breakpoint
CREATE INDEX "income_payment_status_idx" ON "income" USING btree ("payment_status");--> statement-breakpoint
CREATE INDEX "income_reservation_id_idx" ON "income" USING btree ("reservation_id");--> statement-breakpoint
CREATE INDEX "journal_entries_hotel_id_idx" ON "journal_entries" USING btree ("hotel_id");--> statement-breakpoint
CREATE INDEX "journal_entries_entry_number_idx" ON "journal_entries" USING btree ("entry_number");--> statement-breakpoint
CREATE INDEX "journal_entries_entry_date_idx" ON "journal_entries" USING btree ("entry_date");--> statement-breakpoint
CREATE INDEX "journal_entries_status_idx" ON "journal_entries" USING btree ("status");--> statement-breakpoint
CREATE INDEX "journal_lines_entry_id_idx" ON "journal_lines" USING btree ("entry_id");--> statement-breakpoint
CREATE INDEX "journal_lines_account_id_idx" ON "journal_lines" USING btree ("account_id");--> statement-breakpoint
CREATE INDEX "parking_lots_hotel_id_idx" ON "parking_lots" USING btree ("hotel_id");--> statement-breakpoint
CREATE INDEX "parking_lots_status_idx" ON "parking_lots" USING btree ("status");--> statement-breakpoint
CREATE INDEX "parking_spots_lot_id_idx" ON "parking_spots" USING btree ("lot_id");--> statement-breakpoint
CREATE INDEX "parking_spots_status_idx" ON "parking_spots" USING btree ("status");--> statement-breakpoint
CREATE INDEX "parking_spots_spot_number_idx" ON "parking_spots" USING btree ("spot_number");--> statement-breakpoint
CREATE INDEX "parking_transactions_lot_id_idx" ON "parking_transactions" USING btree ("lot_id");--> statement-breakpoint
CREATE INDEX "parking_transactions_guest_id_idx" ON "parking_transactions" USING btree ("guest_id");--> statement-breakpoint
CREATE INDEX "parking_transactions_reservation_id_idx" ON "parking_transactions" USING btree ("reservation_id");--> statement-breakpoint
CREATE INDEX "parking_transactions_status_idx" ON "parking_transactions" USING btree ("status");--> statement-breakpoint
CREATE INDEX "parking_transactions_license_plate_idx" ON "parking_transactions" USING btree ("license_plate");--> statement-breakpoint
CREATE INDEX "reservations_hotel_id_idx" ON "reservations" USING btree ("hotel_id");--> statement-breakpoint
CREATE INDEX "reservations_guest_id_idx" ON "reservations" USING btree ("guest_id");--> statement-breakpoint
CREATE INDEX "reservations_room_id_idx" ON "reservations" USING btree ("room_id");--> statement-breakpoint
CREATE INDEX "reservations_status_idx" ON "reservations" USING btree ("status");--> statement-breakpoint
CREATE INDEX "reservations_check_in_date_idx" ON "reservations" USING btree ("check_in_date");--> statement-breakpoint
CREATE INDEX "reservations_check_out_date_idx" ON "reservations" USING btree ("check_out_date");--> statement-breakpoint
CREATE INDEX "room_types_name_idx" ON "room_types" USING btree ("name");--> statement-breakpoint
CREATE INDEX "rooms_hotel_id_idx" ON "rooms" USING btree ("hotel_id");--> statement-breakpoint
CREATE INDEX "rooms_room_number_idx" ON "rooms" USING btree ("room_number");--> statement-breakpoint
CREATE INDEX "rooms_status_idx" ON "rooms" USING btree ("status");--> statement-breakpoint
CREATE INDEX "rooms_room_type_idx" ON "rooms" USING btree ("room_type");--> statement-breakpoint
CREATE UNIQUE INDEX "user_roles_user_role_idx" ON "user_roles" USING btree ("user_id","role");--> statement-breakpoint
CREATE INDEX "user_roles_user_id_idx" ON "user_roles" USING btree ("user_id");--> statement-breakpoint
CREATE INDEX "users_email_idx" ON "users" USING btree ("email");--> statement-breakpoint
CREATE INDEX "users_role_idx" ON "users" USING btree ("role");--> statement-breakpoint
CREATE INDEX "vehicles_guest_id_idx" ON "vehicles" USING btree ("guest_id");--> statement-breakpoint
CREATE INDEX "vehicles_license_plate_idx" ON "vehicles" USING btree ("license_plate");