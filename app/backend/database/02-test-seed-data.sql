-- SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
--
-- SPDX-License-Identifier: Apache-2.0
--
-- ============================================================================
-- SEED DATA FOR DEVELOPMENT AND TESTING
-- ============================================================================
-- This file contains realistic test data for the Broke da EAR web shop
-- DO NOT RUN IN PRODUCTION
-- ============================================================================
-- PRODUCT CATEGORIES
-- ============================================================================
INSERT INTO
  product_categories (id, name, parent_category_id, sort_order)
VALUES
  (
    '01947f3e-8b2a-7123-b456-123456789001',
    'Synthesizers',
    NULL,
    1
  ),
  (
    '01947f3e-8b2a-7123-b456-123456789002',
    'Effects',
    NULL,
    2
  ),
  (
    '01947f3e-8b2a-7123-b456-123456789003',
    'Utilities',
    NULL,
    3
  ),
  (
    '01947f3e-8b2a-7123-b456-123456789004',
    'Reverb',
    '01947f3e-8b2a-7123-b456-123456789002',
    1
  ),
  (
    '01947f3e-8b2a-7123-b456-123456789005',
    'Delay',
    '01947f3e-8b2a-7123-b456-123456789002',
    2
  ),
  (
    '01947f3e-8b2a-7123-b456-123456789006',
    'Distortion',
    '01947f3e-8b2a-7123-b456-123456789002',
    3
  ),
  (
    '01947f3e-8b2a-7123-b456-123456789007',
    'Virtual Analog',
    '01947f3e-8b2a-7123-b456-123456789001',
    1
  ),
  (
    '01947f3e-8b2a-7123-b456-123456789008',
    'FM Synths',
    '01947f3e-8b2a-7123-b456-123456789001',
    2
  );

-- ============================================================================
-- PRODUCTS
-- ============================================================================
INSERT INTO
  products (
    id,
    name,
    description,
    short_description,
    price,
    currency,
    version,
    download_filename,
    file_size_bytes,
    file_checksum,
    creative_commons_license,
    artistic_credits,
    technical_credits,
    category_id,
    released_at
  )
VALUES
  (
    '01947f3e-8b2a-7123-b456-223456789001',
    'Hawaiian Waves Reverb',
    'Immerse your tracks in the ethereal sound of Hawaiian ocean waves. This reverb plugin captures the natural acoustics of secluded Hawaiian beaches, providing lush, expansive reverbs perfect for ambient music, vocals, and instrumental tracks. Features include adjustable room size, decay time, and a unique "wave modulation" parameter that adds subtle movement to the reverb tail.',
    'Ocean-inspired reverb with Hawaiian beach acoustics and wave modulation.',
    29.99,
    'USD',
    '1.2.0',
    'hawaiian-waves-reverb-v1.2.0.zip',
    45234567,
    'a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456',
    'CC BY-SA 4.0',
    'Field recordings by Koa Nakamura, Maui Sound Collective',
    'DSP Algorithm by Dr. Elena Rodriguez, UI Design by Broke da EAR Team',
    '01947f3e-8b2a-7123-b456-123456789004',
    '2024-11-15 10:00:00-10:00'
  ),
  (
    '01947f3e-8b2a-7123-b456-223456789002',
    'Aloha FM Synthesizer',
    'A powerful FM synthesizer inspired by vintage Hawaiian music and modern electronic sounds. Features 6 operators, multiple algorithm configurations, and a built-in effects chain including chorus, delay, and the signature "Aloha" filter. Perfect for creating everything from classic electric piano sounds to cutting-edge bass lines and lead sounds.',
    'Vintage-inspired FM synthesizer with Hawaiian musical influences.',
    79.99,
    'USD',
    '2.1.3',
    'aloha-fm-synth-v2.1.3.zip',
    128945123,
    'b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef234567',
    NULL,
    'Preset library by Kaleo Johnson, Traditional Hawaiian music research by Aunty Mahina',
    'FM Engine by Broke da EAR Labs, Beta testing by Island Producers Collective',
    '01947f3e-8b2a-7123-b456-123456789008',
    '2024-10-20 14:30:00-10:00'
  ),
  (
    '01947f3e-8b2a-7123-b456-223456789003',
    'Pineapple Compressor',
    'Sweet yet punchy compression inspired by the golden pineapples of Hawaii. This compressor combines the warmth of vintage analog circuits with modern precision control. Features include variable knee, parallel compression blend, and a unique "sweetness" control that adds harmonic coloration. Perfect for drums, vocals, and mix bus processing.',
    'Warm analog-style compressor with harmonic sweetness control.',
    39.99,
    'USD',
    '1.0.5',
    'pineapple-compressor-v1.0.5.zip',
    32156789,
    'c3d4e5f6789012345678901234567890abcdef1234567890abcdef345678',
    NULL,
    NULL,
    'Analog modeling by Dr. Keoni Akamu, Testing by Honolulu Recording Studios',
    '01947f3e-8b2a-7123-b456-123456789002',
    '2024-12-01 09:15:00-10:00'
  ),
  (
    '01947f3e-8b2a-7123-b456-223456789004',
    'Volcano Distortion',
    'Explosive distortion plugin that brings the raw power of Hawaiian volcanoes to your tracks. From subtle tube warmth to devastating fuzz, this plugin offers multiple distortion algorithms inspired by different volcanic formations. Includes pre and post EQ, bias control, and a unique "lava flow" parameter for dynamic saturation.',
    'Powerful distortion with volcanic intensity and multiple algorithms.',
    24.99,
    'USD',
    '1.1.2',
    'volcano-distortion-v1.1.2.zip',
    28934567,
    'd4e5f6789012345678901234567890abcdef1234567890abcdef456789',
    NULL,
    'Inspiration from Big Island volcanic activity recordings',
    'Distortion algorithms by Broke da EAR Research Division',
    '01947f3e-8b2a-7123-b456-123456789006',
    '2024-09-10 16:45:00-10:00'
  ),
  (
    '01947f3e-8b2a-7123-b456-223456789005',
    'Tropical Delay Echo',
    'Multi-tap delay plugin that recreates the natural echo chambers found in Hawaiian lava tubes. Features up to 8 delay taps with individual feedback, filtering, and panning controls. The unique "cave resonance" feature adds natural reverb characteristics to each delay tap, creating complex, evolving echoes perfect for guitars, vocals, and experimental sound design.',
    'Multi-tap delay inspired by Hawaiian lava tube acoustics.',
    34.99,
    'USD',
    '1.3.1',
    'tropical-delay-echo-v1.3.1.zip',
    41267890,
    'e5f6789012345678901234567890abcdef1234567890abcdef567890',
    'CC BY 4.0',
    'Lava tube recordings by Pacific Speleological Society',
    'Echo algorithms by Dr. Leilani Patel, UI/UX by Broke da EAR Design Team',
    '01947f3e-8b2a-7123-b456-123456789005',
    '2024-08-22 11:20:00-10:00'
  ),
  (
    '01947f3e-8b2a-7123-b456-223456789006',
    'Island Utilities Pack',
    'Essential utility plugins for the modern producer. Includes spectrum analyzer with tropical color schemes, tuner with traditional Hawaiian tuning references, and a unique "talk story" communication tool for collaborative sessions. All designed with the laid-back Hawaiian workflow in mind.',
    'Essential utility plugins with Hawaiian-inspired workflow design.',
    19.99,
    'USD',
    '1.0.1',
    'island-utilities-pack-v1.0.1.zip',
    15678923,
    'f6789012345678901234567890abcdef1234567890abcdef678901',
    NULL,
    NULL,
    'Utility design by Broke da EAR Team, Traditional tuning research by Hawaiian Music Institute',
    '01947f3e-8b2a-7123-b456-123456789003',
    '2024-07-04 13:00:00-10:00'
  ),
  (
    '01947f3e-8b2a-7123-b456-223456789007',
    'Beta Synthesizer (Coming Soon)',
    'Next-generation wavetable synthesizer currently in development. Features advanced wavetable morphing, physical modeling elements, and AI-assisted preset generation. This product is not yet available for purchase.',
    'Advanced wavetable synthesizer with AI features (in development).',
    0.00,
    'USD',
    '0.8.0-beta',
    'beta-synth-preview.zip',
    0,
    NULL,
    NULL,
    NULL,
    'Broke da EAR Research & Development Team',
    '01947f3e-8b2a-7123-b456-123456789007',
    '2025-03-01 10:00:00-10:00'
  );

-- ============================================================================
-- USERS
-- ============================================================================
INSERT INTO
  users (
    id,
    auth0_user_id,
    email,
    email_verified,
    password_hash,
    total_purchases_amount,
    total_purchases_count,
    created_at,
    last_login_at
  )
VALUES -- Auth0 users (no password required, placeholder hash used)
  (
    '01947f3e-8b2a-7123-b456-323456789001',
    'auth0|64a7b8c9d0e1f2345678901a',
    'producer@islandbeats.com',
    TRUE,
    'AUTH0_USER_NO_PASSWORD_HASH',
    144.97,
    4,
    '2024-06-15 14:30:00-10:00',
    '2024-12-20 09:45:00-10:00'
  ),
  (
    '01947f3e-8b2a-7123-b456-323456789002',
    'auth0|64a7b8c9d0e1f2345678901b',
    'beats@honolulusound.net',
    TRUE,
    'AUTH0_USER_NO_PASSWORD_HASH',
    79.99,
    1,
    '2024-08-22 16:20:00-10:00',
    '2024-12-19 15:30:00-10:00'
  ),
  -- Email/password users (no Auth0, require password)
  (
    '01947f3e-8b2a-7123-b456-323456789003',
    NULL,
    'music@tropicalstudio.co',
    TRUE,
    '$2a$12$WRx5e3AsDBXJzmh2NS6aHQ0bprAmTYv2yln9RKjIY0L7T8VHG1D0Y',
    59.98,
    2,
    '2024-09-10 11:15:00-10:00',
    '2024-12-18 20:10:00-10:00'
  ),
  (
    '01947f3e-8b2a-7123-b456-323456789004',
    NULL,
    'hello@pacificsounds.org',
    FALSE,
    '$2a$12$XSy6f4BtECYKAni3OT7bIR1cqsBnUZw3zmO0SKkJZ1M8U9WIG2E1Z',
    0.00,
    0,
    '2024-12-01 10:30:00-10:00',
    '2024-12-15 14:20:00-10:00'
  ),
  (
    '01947f3e-8b2a-7123-b456-323456789005',
    NULL,
    'studio@mauivibes.com',
    TRUE,
    '$2a$12$YTz7g5CuFDZLBoj4PU8cJS2drCtObV4AOmP1TLlK02N9V0XJH3F2A',
    29.99,
    1,
    '2024-11-05 13:45:00-10:00',
    '2024-12-22 11:30:00-10:00'
  ),
  -- Mixed example: Auth0 user with unverified email
  (
    '01947f3e-8b2a-7123-b456-323456789006',
    'auth0|64a7b8c9d0e1f2345678901f',
    'sounds@konabeats.io',
    FALSE,
    'AUTH0_USER_NO_PASSWORD_HASH',
    0.00,
    0,
    '2024-12-20 16:45:00-10:00',
    '2024-12-20 16:50:00-10:00'
  );

-- ============================================================================
-- ORDERS
-- ============================================================================
INSERT INTO
  orders (
    id,
    user_id,
    stripe_payment_intent_id,
    stripe_customer_id,
    total_amount,
    currency,
    tax_amount,
    status,
    order_number,
    billing_email,
    billing_name,
    created_at,
    completed_at
  )
VALUES
  (
    '01947f3e-8b2a-7123-b456-423456789001',
    '01947f3e-8b2a-7123-b456-323456789001',
    'pi_3Oj8K12eZvKYlo2C0d7g4f5h',
    'cus_OYxvK8ZQ7P8wXz9A',
    79.99,
    'USD',
    6.80,
    'completed',
    'BDE-20241020-3A7B',
    'producer@islandbeats.com',
    'Kai Nakamura',
    '2024-10-20 15:30:00-10:00',
    '2024-10-20 15:31:45-10:00'
  ),
  (
    '01947f3e-8b2a-7123-b456-423456789002',
    '01947f3e-8b2a-7123-b456-323456789001',
    'pi_3Pk9L23fAvLZmp3D1e8h5g6i',
    'cus_OYxvK8ZQ7P8wXz9A',
    64.98,
    'USD',
    5.52,
    'completed',
    'BDE-20241201-4B8C',
    'producer@islandbeats.com',
    'Kai Nakamura',
    '2024-12-01 10:15:00-10:00',
    '2024-12-01 10:16:22-10:00'
  ),
  (
    '01947f3e-8b2a-7123-b456-423456789003',
    '01947f3e-8b2a-7123-b456-323456789002',
    'pi_3Ql0M34gBwM0nq4E2f9i6h7j',
    'cus_PZywL9AR8Q9xYz0B',
    79.99,
    'USD',
    6.80,
    'completed',
    'BDE-20241110-5C9D',
    'beats@honolulusound.net',
    'Leilani Santos',
    '2024-11-10 14:20:00-10:00',
    '2024-11-10 14:21:18-10:00'
  ),
  (
    '01947f3e-8b2a-7123-b456-423456789004',
    '01947f3e-8b2a-7123-b456-323456789003',
    'pi_3Rm1N45hCxN1or5F3g0j7i8k',
    'cus_QAzxM0BS9R0yZz1C',
    59.98,
    'USD',
    5.10,
    'completed',
    'BDE-20241115-6D0E',
    'music@tropicalstudio.co',
    'Diego Rodriguez',
    '2024-11-15 16:45:00-10:00',
    '2024-11-15 16:46:35-10:00'
  ),
  (
    '01947f3e-8b2a-7123-b456-423456789005',
    '01947f3e-8b2a-7123-b456-323456789005',
    'pi_3Sn2O56iDyO2ps6G4h1k8j9l',
    'cus_RBAyN1CT0S1zA02D',
    29.99,
    'USD',
    2.55,
    'completed',
    'BDE-20241205-7E1F',
    'studio@mauivibes.com',
    'Makoa Johnson',
    '2024-12-05 12:30:00-10:00',
    '2024-12-05 12:31:12-10:00'
  ),
  (
    '01947f3e-8b2a-7123-b456-423456789006',
    '01947f3e-8b2a-7123-b456-323456789004',
    'pi_3To3P67jEzP3qt7H5i2l9k0m',
    'cus_SCBzO2DU1T2zB13E',
    39.99,
    'USD',
    3.40,
    'pending',
    'BDE-20241222-8F2G',
    'hello@pacificsounds.org',
    'Emma Chen',
    '2024-12-22 09:15:00-10:00',
    NULL
  );

-- ============================================================================
-- ORDER ITEMS
-- ============================================================================
INSERT INTO
  order_items (
    id,
    order_id,
    product_id,
    product_name,
    product_price,
    quantity,
    line_total,
    created_at
  )
VALUES
  (
    '01947f3e-8b2a-7123-b456-523456789001',
    '01947f3e-8b2a-7123-b456-423456789001',
    '01947f3e-8b2a-7123-b456-223456789002',
    'Aloha FM Synthesizer',
    79.99,
    1,
    79.99,
    '2024-10-20 15:30:00-10:00'
  ),
  (
    '01947f3e-8b2a-7123-b456-523456789002',
    '01947f3e-8b2a-7123-b456-423456789002',
    '01947f3e-8b2a-7123-b456-223456789003',
    'Pineapple Compressor',
    39.99,
    1,
    39.99,
    '2024-12-01 10:15:00-10:00'
  ),
  (
    '01947f3e-8b2a-7123-b456-523456789003',
    '01947f3e-8b2a-7123-b456-423456789002',
    '01947f3e-8b2a-7123-b456-223456789004',
    'Volcano Distortion',
    24.99,
    1,
    24.99,
    '2024-12-01 10:15:00-10:00'
  ),
  (
    '01947f3e-8b2a-7123-b456-523456789004',
    '01947f3e-8b2a-7123-b456-423456789003',
    '01947f3e-8b2a-7123-b456-223456789002',
    'Aloha FM Synthesizer',
    79.99,
    1,
    79.99,
    '2024-11-10 14:20:00-10:00'
  ),
  (
    '01947f3e-8b2a-7123-b456-523456789005',
    '01947f3e-8b2a-7123-b456-423456789004',
    '01947f3e-8b2a-7123-b456-223456789001',
    'Hawaiian Waves Reverb',
    29.99,
    1,
    29.99,
    '2024-11-15 16:45:00-10:00'
  ),
  (
    '01947f3e-8b2a-7123-b456-523456789006',
    '01947f3e-8b2a-7123-b456-423456789004',
    '01947f3e-8b2a-7123-b456-223456789001',
    'Hawaiian Waves Reverb',
    29.99,
    1,
    29.99,
    '2024-11-15 16:45:00-10:00'
  ),
  (
    '01947f3e-8b2a-7123-b456-523456789007',
    '01947f3e-8b2a-7123-b456-423456789005',
    '01947f3e-8b2a-7123-b456-223456789001',
    'Hawaiian Waves Reverb',
    29.99,
    1,
    29.99,
    '2024-12-05 12:30:00-10:00'
  ),
  (
    '01947f3e-8b2a-7123-b456-523456789008',
    '01947f3e-8b2a-7123-b456-423456789006',
    '01947f3e-8b2a-7123-b456-223456789003',
    'Pineapple Compressor',
    39.99,
    1,
    39.99,
    '2024-12-22 09:15:00-10:00'
  );

-- ============================================================================
-- USER DOWNLOADS
-- ============================================================================
INSERT INTO
  user_downloads (
    id,
    user_id,
    product_id,
    order_id,
    download_count,
    last_downloaded_at,
    created_at
  )
VALUES
  (
    '01947f3e-8b2a-7123-b456-623456789001',
    '01947f3e-8b2a-7123-b456-323456789001',
    '01947f3e-8b2a-7123-b456-223456789002',
    '01947f3e-8b2a-7123-b456-423456789001',
    3,
    '2024-12-15 10:30:00-10:00',
    '2024-10-20 15:32:00-10:00'
  ),
  (
    '01947f3e-8b2a-7123-b456-623456789002',
    '01947f3e-8b2a-7123-b456-323456789001',
    '01947f3e-8b2a-7123-b456-223456789003',
    '01947f3e-8b2a-7123-b456-423456789002',
    1,
    '2024-12-01 10:17:00-10:00',
    '2024-12-01 10:17:00-10:00'
  ),
  (
    '01947f3e-8b2a-7123-b456-623456789003',
    '01947f3e-8b2a-7123-b456-323456789001',
    '01947f3e-8b2a-7123-b456-223456789004',
    '01947f3e-8b2a-7123-b456-423456789002',
    2,
    '2024-12-10 14:20:00-10:00',
    '2024-12-01 10:17:00-10:00'
  ),
  (
    '01947f3e-8b2a-7123-b456-623456789004',
    '01947f3e-8b2a-7123-b456-323456789002',
    '01947f3e-8b2a-7123-b456-223456789002',
    '01947f3e-8b2a-7123-b456-423456789003',
    1,
    '2024-11-10 14:22:00-10:00',
    '2024-11-10 14:22:00-10:00'
  ),
  (
    '01947f3e-8b2a-7123-b456-623456789005',
    '01947f3e-8b2a-7123-b456-323456789003',
    '01947f3e-8b2a-7123-b456-223456789001',
    '01947f3e-8b2a-7123-b456-423456789004',
    2,
    '2024-12-01 09:45:00-10:00',
    '2024-11-15 16:47:00-10:00'
  ),
  (
    '01947f3e-8b2a-7123-b456-623456789006',
    '01947f3e-8b2a-7123-b456-323456789005',
    '01947f3e-8b2a-7123-b456-223456789001',
    '01947f3e-8b2a-7123-b456-423456789005',
    1,
    '2024-12-05 12:32:00-10:00',
    '2024-12-05 12:32:00-10:00'
  );
