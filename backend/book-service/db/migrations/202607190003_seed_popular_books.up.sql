ALTER TABLE books
	ADD COLUMN cover_url VARCHAR(512) NULL AFTER description;

-- Snapshot of Open Library's monthly trending catalog retrieved on 2026-07-19.
-- INSERT IGNORE preserves any catalog records already registered with the same ISBN.
INSERT IGNORE INTO books (
	id, isbn, title, author, description, cover_url, publication_year,
	total_copies, available_copies, created_at, updated_at
) VALUES
(
	'10000000-0000-4000-8000-000000000001',
	'9781847941831',
	'Atomic Habits',
	'James Clear',
	'A practical guide to forming good habits, breaking bad ones, and mastering the small behaviors that lead to remarkable results.',
	'https://covers.openlibrary.org/b/id/12539702-M.jpg',
	2016, 5, 5, '2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
) ,
(
	'10000000-0000-4000-8000-000000000002',
	'9780140280197',
	'The 48 Laws of Power',
	'Robert Greene',
	'A study of power that distills historical examples and strategic lessons into forty-eight laws.',
	'https://covers.openlibrary.org/b/id/6424160-M.jpg',
	1998, 5, 5, '2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
) ,
(
	'10000000-0000-4000-8000-000000000003',
	'9780446567404',
	'Rich Dad, Poor Dad',
	'Robert T. Kiyosaki, Sharon L. Lechter',
	'A personal-finance classic about building financial literacy, assets, and long-term independence.',
	'https://covers.openlibrary.org/b/id/8315603-M.jpg',
	1990, 5, 5, '2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
) ,
(
	'10000000-0000-4000-8000-000000000004',
	'9781408855898',
	'Harry Potter and the Philosopher''s Stone',
	'J. K. Rowling',
	'Harry Potter discovers that he is a wizard and begins an extraordinary adventure at Hogwarts School of Witchcraft and Wizardry.',
	'https://covers.openlibrary.org/b/id/15155833-M.jpg',
	1997, 5, 5, '2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
) ,
(
	'10000000-0000-4000-8000-000000000005',
	'9781804090114',
	'The Psychology of Money',
	'Morgan Housel',
	'Nineteen stories exploring how behavior, personal history, and incentives shape the way people think about money.',
	'https://covers.openlibrary.org/b/id/10389354-M.jpg',
	2020, 5, 5, '2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
) ,
(
	'10000000-0000-4000-8000-000000000006',
	'9780671027032',
	'How to Win Friends and Influence People',
	'Dale Carnegie',
	'A practical guide to communicating effectively, building relationships, and leading with empathy.',
	'https://covers.openlibrary.org/b/id/13314878-M.jpg',
	1936, 5, 5, '2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
) ,
(
	'10000000-0000-4000-8000-000000000007',
	'9781508706250',
	'The Deal',
	'Elle Kennedy',
	'A college sports romance in which a tutoring arrangement and a pretend date develop into something real.',
	'https://covers.openlibrary.org/b/id/10201611-M.jpg',
	2000, 5, 5, '2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
) ,
(
	'10000000-0000-4000-8000-000000000008',
	'9780143110163',
	'Think and Grow Rich',
	'Napoleon Hill',
	'A classic examination of the principles, mindset, and habits associated with professional success.',
	'https://covers.openlibrary.org/b/id/14542536-M.jpg',
	1937, 5, 5, '2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
) ,
(
	'10000000-0000-4000-8000-000000000009',
	'0671663984',
	'The 7 Habits of Highly Effective People',
	'Stephen R. Covey, Sean Covey',
	'A principle-centered framework for personal effectiveness, collaboration, and continuous improvement.',
	'https://covers.openlibrary.org/b/id/10079937-M.jpg',
	1989, 5, 5, '2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
) ,
(
	'10000000-0000-4000-8000-000000000010',
	'0061122416',
	'The Alchemist',
	'Paulo Coelho',
	'A young shepherd travels in search of treasure and discovers lessons about purpose, courage, and following his dreams.',
	'https://covers.openlibrary.org/b/id/7414780-M.jpg',
	1988, 5, 5, '2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
);