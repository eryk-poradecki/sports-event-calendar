-- Venues
INSERT INTO venues (name, city, _country_id, address, capacity, website_url, created_at, updated_at)
VALUES
    (
        'Ernst-Happel-Stadion',
        'Vienna',
        (SELECT id FROM countries WHERE code_alpha2 = 'AT'),
        'Meiereistraße 7, 1020 Vienna',
        50865,
        'https://en.wikipedia.org/wiki/Ernst-Happel-Stadion',
        now(),
        now()
    ),
    (
        'Red Bull Arena',
        'Salzburg',
        (SELECT id FROM countries WHERE code_alpha2 = 'AT'),
        'Stadionstraße 2/3, 5071 Wals-Siezenheim',
        30188,
        'https://www.redbullsalzburg.at/en/red-bull-arena/',
        now(),
        now()
    ),
    (
        'Wiener Stadthalle',
        'Vienna',
        (SELECT id FROM countries WHERE code_alpha2 = 'AT'),
        'Roland-Rainer-Platz 1, 1150 Vienna',
        16000,
        'https://www.stadthalle.com/en',
        now(),
        now()
    ),
    (
        'Stadthalle Graz',
        'Graz',
        (SELECT id FROM countries WHERE code_alpha2 = 'AT'),
        'Messeplatz 1, 8010 Graz',
        11500,
        'https://en.wikipedia.org/wiki/Stadthalle_Graz',
        now(),
        now()
    )
ON CONFLICT DO NOTHING;

-- Competitions
INSERT INTO competitions (name, slug, _sport_id, start_date, end_date, description, created_at, updated_at)
VALUES
    (
        'Austrian Bundesliga',
        'austrian-bundesliga',
        (SELECT id FROM sports WHERE slug = 'football'),
        DATE '2025-07-01',
        DATE '2026-05-31',
        'Top-tier football league in Austria.',
        now(),
        now()
    ),
    (
        'ICE Hockey League',
        'ice-hockey-league',
        (SELECT id FROM sports WHERE slug = 'ice-hockey'),
        DATE '2025-09-01',
        DATE '2026-04-30',
        'Professional ice hockey league in Central Europe.',
        now(),
        now()
    ),
    (
        'Austrian Basketball Superliga',
        'austrian-basketball-superliga',
        (SELECT id FROM sports WHERE slug = 'basketball'),
        DATE '2025-09-01',
        DATE '2026-06-15',
        'Top-level Austrian basketball competition.',
        now(),
        now()
    ),
    (
        'Austrian Volley League Men',
        'austrian-volley-league-men',
        (SELECT id FROM sports WHERE slug = 'volleyball'),
        DATE '2025-10-01',
        DATE '2026-05-15',
        'Top-level men volleyball competition in Austria.',
        now(),
        now()
    )
ON CONFLICT (slug) DO NOTHING;

-- Teams
INSERT INTO teams (name, slug, _country_id, _sport_id, website_url, created_at, updated_at)
VALUES
    (
        'Red Bull Salzburg',
        'red-bull-salzburg-football',
        (SELECT id FROM countries WHERE code_alpha2 = 'AT'),
        (SELECT id FROM sports WHERE slug = 'football'),
        'https://www.redbullsalzburg.at/en',
        now(),
        now()
    ),
    (
        'SK Sturm Graz',
        'sk-sturm-graz-football',
        (SELECT id FROM countries WHERE code_alpha2 = 'AT'),
        (SELECT id FROM sports WHERE slug = 'football'),
        'https://sksturm.at/',
        now(),
        now()
    ),
    (
        'Rapid Wien',
        'rapid-wien-football',
        (SELECT id FROM countries WHERE code_alpha2 = 'AT'),
        (SELECT id FROM sports WHERE slug = 'football'),
        'https://www.skrapid.at/',
        now(),
        now()
    ),
    (
        'LASK',
        'lask-football',
        (SELECT id FROM countries WHERE code_alpha2 = 'AT'),
        (SELECT id FROM sports WHERE slug = 'football'),
        'https://www.lask.at/en/m',
        now(),
        now()
    ),
    (
        'EC KAC',
        'ec-kac-ice-hockey',
        (SELECT id FROM countries WHERE code_alpha2 = 'AT'),
        (SELECT id FROM sports WHERE slug = 'ice-hockey'),
        'https://www.kac.at/',
        now(),
        now()
    ),
    (
        'Vienna Capitals',
        'vienna-capitals-ice-hockey',
        (SELECT id FROM countries WHERE code_alpha2 = 'AT'),
        (SELECT id FROM sports WHERE slug = 'ice-hockey'),
        'https://www.vienna-capitals.at/',
        now(),
        now()
    ),
    (
        'BC Vienna',
        'bc-vienna-basketball',
        (SELECT id FROM countries WHERE code_alpha2 = 'AT'),
        (SELECT id FROM sports WHERE slug = 'basketball'),
        'https://bcvienna.com/',
        now(),
        now()
    ),
    (
        'Oberwart Gunners',
        'oberwart-gunners-basketball',
        (SELECT id FROM countries WHERE code_alpha2 = 'AT'),
        (SELECT id FROM sports WHERE slug = 'basketball'),
        'https://www.gunners.at/',
        now(),
        now()
    ),
    (
        'Hypo Tirol Volleyballteam',
        'hypo-tirol-volleyball',
        (SELECT id FROM countries WHERE code_alpha2 = 'AT'),
        (SELECT id FROM sports WHERE slug = 'volleyball'),
        'https://volleyballteamtirol.com/',
        now(),
        now()
    ),
    (
        'UVC Graz',
        'uvc-graz-volleyball',
        (SELECT id FROM countries WHERE code_alpha2 = 'AT'),
        (SELECT id FROM sports WHERE slug = 'volleyball'),
        'https://www.uvcgraz.at/',
        now(),
        now()
    )
ON CONFLICT (slug) DO NOTHING;

-- Events
INSERT INTO events (
    _sport_id,
    _competition_id,
    _venue_id,
    _home_team_id,
    _away_team_id,
    start_time,
    status,
    home_score,
    away_score,
    description,
    is_neutral_venue,
    created_at,
    updated_at
)
VALUES
    (
        (SELECT id FROM sports WHERE slug = 'football'),
        (SELECT id FROM competitions WHERE slug = 'austrian-bundesliga'),
        (SELECT id FROM venues WHERE name = 'Red Bull Arena'),
        (SELECT id FROM teams WHERE slug = 'red-bull-salzburg-football'),
        (SELECT id FROM teams WHERE slug = 'sk-sturm-graz-football'),
        TIMESTAMPTZ '2026-03-22 17:00:00+01',
        'scheduled',
        NULL,
        NULL,
        'Bundesliga matchday fixture.',
        FALSE,
        now(),
        now()
    ),
    (
        (SELECT id FROM sports WHERE slug = 'football'),
        (SELECT id FROM competitions WHERE slug = 'austrian-bundesliga'),
        (SELECT id FROM venues WHERE name = 'Ernst-Happel-Stadion'),
        (SELECT id FROM teams WHERE slug = 'rapid-wien-football'),
        (SELECT id FROM teams WHERE slug = 'lask-football'),
        TIMESTAMPTZ '2026-03-23 19:30:00+01',
        'scheduled',
        NULL,
        NULL,
        'Evening league fixture in Vienna.',
        FALSE,
        now(),
        now()
    ),
    (
        (SELECT id FROM sports WHERE slug = 'football'),
        (SELECT id FROM competitions WHERE slug = 'austrian-bundesliga'),
        (SELECT id FROM venues WHERE name = 'Red Bull Arena'),
        (SELECT id FROM teams WHERE slug = 'red-bull-salzburg-football'),
        (SELECT id FROM teams WHERE slug = 'lask-football'),
        TIMESTAMPTZ '2026-03-10 18:30:00+01',
        'finished',
        3,
        1,
        'Finished league fixture.',
        FALSE,
        now(),
        now()
    ),
    (
        (SELECT id FROM sports WHERE slug = 'football'),
        (SELECT id FROM competitions WHERE slug = 'austrian-bundesliga'),
        (SELECT id FROM venues WHERE name = 'Ernst-Happel-Stadion'),
        (SELECT id FROM teams WHERE slug = 'rapid-wien-football'),
        (SELECT id FROM teams WHERE slug = 'sk-sturm-graz-football'),
        TIMESTAMPTZ '2026-03-05 20:00:00+01',
        'finished',
        2,
        2,
        'Draw after a late equalizer.',
        FALSE,
        now(),
        now()
    ),
    (
        (SELECT id FROM sports WHERE slug = 'football'),
        (SELECT id FROM competitions WHERE slug = 'austrian-bundesliga'),
        (SELECT id FROM venues WHERE name = 'Ernst-Happel-Stadion'),
        (SELECT id FROM teams WHERE slug = 'rapid-wien-football'),
        (SELECT id FROM teams WHERE slug = 'red-bull-salzburg-football'),
        TIMESTAMPTZ '2026-04-02 18:00:00+01',
        'cancelled',
        NULL,
        NULL,
        'Cancelled due to severe weather conditions.',
        FALSE,
        now(),
        now()
    ),

    (
        (SELECT id FROM sports WHERE slug = 'ice-hockey'),
        (SELECT id FROM competitions WHERE slug = 'ice-hockey-league'),
        (SELECT id FROM venues WHERE name = 'Wiener Stadthalle'),
        (SELECT id FROM teams WHERE slug = 'vienna-capitals-ice-hockey'),
        (SELECT id FROM teams WHERE slug = 'ec-kac-ice-hockey'),
        TIMESTAMPTZ '2026-03-24 19:15:00+01',
        'scheduled',
        NULL,
        NULL,
        'Regular season hockey match.',
        FALSE,
        now(),
        now()
    ),
    (
        (SELECT id FROM sports WHERE slug = 'ice-hockey'),
        (SELECT id FROM competitions WHERE slug = 'ice-hockey-league'),
        (SELECT id FROM venues WHERE name = 'Wiener Stadthalle'),
        (SELECT id FROM teams WHERE slug = 'vienna-capitals-ice-hockey'),
        (SELECT id FROM teams WHERE slug = 'ec-kac-ice-hockey'),
        TIMESTAMPTZ '2026-03-08 18:45:00+01',
        'finished',
        1,
        4,
        'Strong away performance by KAC.',
        FALSE,
        now(),
        now()
    ),

    (
        (SELECT id FROM sports WHERE slug = 'basketball'),
        (SELECT id FROM competitions WHERE slug = 'austrian-basketball-superliga'),
        (SELECT id FROM venues WHERE name = 'Wiener Stadthalle'),
        (SELECT id FROM teams WHERE slug = 'bc-vienna-basketball'),
        (SELECT id FROM teams WHERE slug = 'oberwart-gunners-basketball'),
        TIMESTAMPTZ '2026-03-26 18:00:00+01',
        'scheduled',
        NULL,
        NULL,
        'Top-flight basketball fixture.',
        FALSE,
        now(),
        now()
    ),
    (
        (SELECT id FROM sports WHERE slug = 'basketball'),
        (SELECT id FROM competitions WHERE slug = 'austrian-basketball-superliga'),
        (SELECT id FROM venues WHERE name = 'Wiener Stadthalle'),
        (SELECT id FROM teams WHERE slug = 'bc-vienna-basketball'),
        (SELECT id FROM teams WHERE slug = 'oberwart-gunners-basketball'),
        TIMESTAMPTZ '2026-03-12 18:00:00+01',
        'finished',
        88,
        81,
        'BC Vienna closed the game strongly in the fourth quarter.',
        FALSE,
        now(),
        now()
    ),

    (
        (SELECT id FROM sports WHERE slug = 'volleyball'),
        (SELECT id FROM competitions WHERE slug = 'austrian-volley-league-men'),
        (SELECT id FROM venues WHERE name = 'Stadthalle Graz'),
        (SELECT id FROM teams WHERE slug = 'uvc-graz-volleyball'),
        (SELECT id FROM teams WHERE slug = 'hypo-tirol-volleyball'),
        TIMESTAMPTZ '2026-03-27 19:00:00+01',
        'scheduled',
        NULL,
        NULL,
        'League volleyball fixture.',
        FALSE,
        now(),
        now()
    ),
    (
        (SELECT id FROM sports WHERE slug = 'volleyball'),
        (SELECT id FROM competitions WHERE slug = 'austrian-volley-league-men'),
        (SELECT id FROM venues WHERE name = 'Stadthalle Graz'),
        (SELECT id FROM teams WHERE slug = 'uvc-graz-volleyball'),
        (SELECT id FROM teams WHERE slug = 'hypo-tirol-volleyball'),
        TIMESTAMPTZ '2026-03-14 19:00:00+01',
        'finished',
        2,
        3,
        'Five-set thriller won by Hypo Tirol.',
        FALSE,
        now(),
        now()
    ),

    (
        (SELECT id FROM sports WHERE slug = 'football'),
        NULL,
        NULL,
        (SELECT id FROM teams WHERE slug = 'sk-sturm-graz-football'),
        (SELECT id FROM teams WHERE slug = 'lask-football'),
        TIMESTAMPTZ '2026-04-10 20:15:00+01',
        'scheduled',
        NULL,
        NULL,
        'Example event without competition and venue.',
        FALSE,
        now(),
        now()
    ),
    (
        (SELECT id FROM sports WHERE slug = 'football'),
        (SELECT id FROM competitions WHERE slug = 'austrian-bundesliga'),
        (SELECT id FROM venues WHERE name = 'Red Bull Arena'),
        (SELECT id FROM teams WHERE slug = 'red-bull-salzburg-football'),
        (SELECT id FROM teams WHERE slug = 'rapid-wien-football'),
        TIMESTAMPTZ '2026-04-14 20:30:00+01',
        'scheduled',
        NULL,
        NULL,
        'Late-season headline fixture.',
        FALSE,
        now(),
        now()
    );