-- Snapshot of Open Library's monthly trending catalog retrieved on 2026-07-19.
-- INSERT IGNORE preserves catalog records already registered with the same ISBN.
INSERT IGNORE INTO books (
	id, isbn, title, author, description, cover_url, publication_year,
	total_copies, available_copies, created_at, updated_at
) VALUES
(
	'10000000-0000-4000-8000-000000000011', '9781407135397', 'The Hunger Games', 'Suzanne Collins', NULL,
	'https://covers.openlibrary.org/b/id/12646537-M.jpg', 2008, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000012', '9780525428145', 'The Laws of Human Nature', 'Robert Greene', NULL,
	'https://covers.openlibrary.org/b/id/10170095-M.jpg', 2018, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000013', '9780385676519', 'Thinking, fast and slow', 'Daniel Kahneman', NULL,
	'https://covers.openlibrary.org/b/id/13290711-M.jpg', 2011, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000014', '0198185219', 'Nineteen Eighty-Four', 'George Orwell', NULL,
	'https://covers.openlibrary.org/b/id/9267242-M.jpg', 1949, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000015', '0894643991', 'The psychology of money', 'Henry Clay Lindgren', NULL,
	'https://covers.openlibrary.org/b/id/6694177-M.jpg', 1991, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000016', '9780062899149', 'The Subtle Art of Not Giving a F*ck', 'Mark Manson', NULL,
	'https://covers.openlibrary.org/b/id/8231990-M.jpg', 2016, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000017', '9781405664387', 'The Lightning Thief', 'Rick Riordan', NULL,
	'https://covers.openlibrary.org/b/id/7239831-M.jpg', 2005, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000018', '9780385539258', 'A Little Life', 'Hanya Yanagihara', NULL,
	'https://covers.openlibrary.org/b/id/12065783-M.jpg', 2008, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000019', '9781524796280', 'Fire & Blood', 'George R. R. Martin', NULL,
	'https://covers.openlibrary.org/b/id/12063529-M.jpg', 2014, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000020', '0810993139', 'Diary of a Wimpy Kid', 'Jeff Kinney', NULL,
	'https://covers.openlibrary.org/b/id/14376136-M.jpg', 2007, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000021', '9780062085481', 'Shatter Me', 'Tahereh Mafi', NULL,
	'https://covers.openlibrary.org/b/id/6974992-M.jpg', 2011, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000022', '8459912019', 'Le petit prince', 'Antoine de Saint-Exupéry', NULL,
	'https://covers.openlibrary.org/b/id/10708272-M.jpg', 1943, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000023', '9798217116775', 'A Good Girl''s Guide to Murder', 'Holly Jackson', NULL,
	'https://covers.openlibrary.org/b/id/13156188-M.jpg', 2019, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000024', '1861977697', 'The Art of Seduction', 'Robert Greene, Joost Elffers', NULL,
	'https://covers.openlibrary.org/b/id/917840-M.jpg', 2001, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000025', '9788192910963', 'The Power of Your Subconscious Mind', 'Joseph Murphy', NULL,
	'https://covers.openlibrary.org/b/id/6553019-M.jpg', 1963, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000026', '9781501197277', '嫌われる勇気', 'Ichirō Kishimi, Fumitake Koga', NULL,
	'https://covers.openlibrary.org/b/id/10873626-M.jpg', 2013, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000027', '9780593395561', 'Project Hail Mary', 'Andy Weir', NULL,
	'https://covers.openlibrary.org/b/id/11200092-M.jpg', 2021, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000028', '9780316015844', 'Twilight', 'Stephenie Meyer', NULL,
	'https://covers.openlibrary.org/b/id/12641977-M.jpg', 2005, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000029', '9781471158254', 'It Ends With Us', 'Colleen Hoover', NULL,
	'https://covers.openlibrary.org/b/id/10473609-M.jpg', 2012, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000030', '9781949759228', 'The Mountain Is You', 'Brianna Wiest', NULL,
	'https://covers.openlibrary.org/b/id/13838236-M.jpg', 2020, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000031', '9781537356426', 'The Mistake', 'Elle Kennedy', NULL,
	'https://covers.openlibrary.org/b/id/10420220-M.jpg', 2015, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000032', '0747538484', 'Harry Potter and the Chamber of Secrets', 'J. K. Rowling', NULL,
	'https://covers.openlibrary.org/b/id/15158664-M.jpg', 1998, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000033', '9781455586691', 'Deep Work', 'Cal Newport', NULL,
	'https://covers.openlibrary.org/b/id/7988607-M.jpg', 2016, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000034', '9780771038501', 'Sapiens', 'Yuval Noah Harari', NULL,
	'https://covers.openlibrary.org/b/id/8634250-M.jpg', 2011, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000035', '9798395421142', 'The Richest Man in Babylon', 'George S. Clason', NULL,
	'https://covers.openlibrary.org/b/id/10491331-M.jpg', 1926, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000036', '072161566X', 'Psychology', 'Robert A. Baron', NULL,
	'https://covers.openlibrary.org/b/id/1146622-M.jpg', 1977, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000037', '9780143130727', 'Ikigai', 'Héctor García, Francesc Miralles', NULL,
	'https://covers.openlibrary.org/b/id/11300391-M.jpg', 2016, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000038', '9781780332000', 'Mindset', 'Carol S. Dweck', NULL,
	'https://covers.openlibrary.org/b/id/746414-M.jpg', 2006, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000039', '9781408855911', 'Harry Potter and the Prisoner of Azkaban', 'J. K. Rowling', NULL,
	'https://covers.openlibrary.org/b/id/10580435-M.jpg', 1999, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000040', '0395362903', 'The Hobbit', 'J.R.R. Tolkien', NULL,
	'https://covers.openlibrary.org/b/id/14627509-M.jpg', 1703, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000041', '0451149513', 'It', 'Stephen King', NULL,
	'https://covers.openlibrary.org/b/id/8569284-M.jpg', 1986, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000042', '9781408806159', 'The Kite Runner', 'Khaled Hosseini', NULL,
	'https://covers.openlibrary.org/b/id/14846827-M.jpg', 2003, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000043', '007141858X', 'How to Talk to Anyone', 'Leil Lowndes', NULL,
	'https://covers.openlibrary.org/b/id/58950-M.jpg', 1999, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000044', '9780735201446', 'Trading in the zone', 'Mark Douglas', NULL,
	'https://covers.openlibrary.org/b/id/460625-M.jpg', 2001, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000045', '9780062442390', 'Red Queen', 'Victoria Aveyard', NULL,
	'https://covers.openlibrary.org/b/id/7434883-M.jpg', 2015, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000046', '9781405664417', 'The Battle of the Labyrinth', 'Rick Riordan', NULL,
	'https://covers.openlibrary.org/b/id/6274739-M.jpg', 2005, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000047', '9780141321288', 'The Last Olympian', 'Rick Riordan', NULL,
	'https://covers.openlibrary.org/b/id/6624107-M.jpg', 2009, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000048', '9780345539786', 'Red Rising', 'Pierce Brown', NULL,
	'https://covers.openlibrary.org/b/id/7316188-M.jpg', 2014, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000049', '9781500420536', 'Suicide Med', 'Freida McFadden', NULL,
	'https://covers.openlibrary.org/b/id/13470475-M.jpg', 2014, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000050', '0446693219', 'Rich Dad, Poor Dad for Teens', 'Robert T. Kiyosaki, Sharon L. Lechter', NULL,
	'https://covers.openlibrary.org/b/id/1174758-M.jpg', 2004, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000051', '9781406365931', 'The Hired Girl', 'Laura Amy Schlitz', NULL,
	'https://covers.openlibrary.org/b/id/11610624-M.jpg', 2015, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000052', '9780078038198', 'The Power of Logic', 'Frances Howard-Snyder, Daniel Howard-Snyder, Ryan Wasserman', NULL,
	'https://covers.openlibrary.org/b/id/7157258-M.jpg', 2009, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000053', '9781423113461', 'The Lost Hero', 'Rick Riordan', NULL,
	'https://covers.openlibrary.org/b/id/12848687-M.jpg', 2010, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000054', '9780525478812', 'The Fault in Our Stars', 'John Green', NULL,
	'https://covers.openlibrary.org/b/id/7418786-M.jpg', 2010, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000055', '9781737475736', '$100M Offers', 'Alex Hormozi', NULL,
	'https://covers.openlibrary.org/b/id/11948182-M.jpg', 2021, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000056', '9781532902949', 'Brain Damage', 'Freida McFadden', NULL,
	'https://covers.openlibrary.org/b/id/15132892-M.jpg', 2016, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000057', '9781804955871', 'God of Malice', 'Rina Kent', NULL,
	'https://covers.openlibrary.org/b/id/13291765-M.jpg', 2022, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000058', '0312146027', 'Forbidden knowledge', 'Roger Shattuck', NULL,
	'https://covers.openlibrary.org/b/id/8295312-M.jpg', 1996, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000059', '9781407405445', 'A Thousand Splendid Suns', 'Khaled Hosseini', NULL,
	'https://covers.openlibrary.org/b/id/8579790-M.jpg', 2007, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000060', '9781518783876', 'Corrupt', 'Penelope Douglas', NULL,
	'https://covers.openlibrary.org/b/id/10226443-M.jpg', 2015, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000061', '0688015603', 'Influence', 'Robert B. Cialdini', NULL,
	'https://covers.openlibrary.org/b/id/431011-M.jpg', 1983, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000062', '9781444741292', 'Misery', 'Stephen King', NULL,
	'https://covers.openlibrary.org/b/id/8259296-M.jpg', 1978, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000063', '9781537356730', 'The Score', 'Elle Kennedy', NULL,
	'https://covers.openlibrary.org/b/id/10538493-M.jpg', 2016, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000064', '9781524700805', 'Girl in Pieces', 'Kathleen Glasgow', NULL,
	'https://covers.openlibrary.org/b/id/8888850-M.jpg', 2000, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000065', '9781250179937', 'Surrounded by Idiots', 'Thomas Erikson', NULL,
	'https://covers.openlibrary.org/b/id/10105591-M.jpg', 2014, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000066', '9781635575552', 'A Court of Thorns and Roses', 'Sarah J. Maas', NULL,
	'https://covers.openlibrary.org/b/id/8738585-M.jpg', 2013, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000067', '9780141321264', 'The Titan''s Curse', 'Rick Riordan', NULL,
	'https://covers.openlibrary.org/b/id/14601475-M.jpg', 2007, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000068', '9780375869020', 'Wonder', 'R. J. Palacio', NULL,
	'https://covers.openlibrary.org/b/id/8223160-M.jpg', 2001, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000069', '0140278761', 'Cien años de soledad', 'Gabriel García Márquez', NULL,
	'https://covers.openlibrary.org/b/id/12627383-M.jpg', 1967, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000070', '0670785938', 'The Body Keeps the Score', 'Bessel van der Kolk', NULL,
	'https://covers.openlibrary.org/b/id/8315367-M.jpg', 2014, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000071', '9789382563792', 'Emotional Intelligence', 'Daniel Goleman', NULL,
	'https://covers.openlibrary.org/b/id/1359485-M.jpg', 1995, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000072', '9781484160916', 'Jane Eyre', 'Charlotte Brontë', NULL,
	'https://covers.openlibrary.org/b/id/8235363-M.jpg', 1847, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000073', '9781495407956', 'The 48 Laws of Power Pivotal Points -The Pivotal Guide to Robert Greene''s Celebrated Book', 'Pivotal Point Papers', NULL,
	'https://covers.openlibrary.org/b/id/14437046-M.jpg', 2014, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000074', '9781582701707', 'The Secret', 'Rhonda Byrne', NULL,
	'https://covers.openlibrary.org/b/id/845815-M.jpg', 2000, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000075', '9780329143558', 'I, Robot', 'Isaac Asimov', NULL,
	'https://covers.openlibrary.org/b/id/12385229-M.jpg', 1950, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000076', '9780545010221', 'Harry Potter and the Deathly Hallows', 'J. K. Rowling', NULL,
	'https://covers.openlibrary.org/b/id/15158660-M.jpg', 2007, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000077', '9780141385297', 'Diary of a Wimpy Kid', 'Jeff Kinney', NULL,
	'https://covers.openlibrary.org/b/id/10332154-M.jpg', 2017, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000078', '9781405664394', 'The Sea of Monsters', 'Rick Riordan', NULL,
	'https://covers.openlibrary.org/b/id/108909-M.jpg', 2005, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000079', '9780007659579', 'Murder on the Orient Express', 'Agatha Christie', NULL,
	'https://covers.openlibrary.org/b/id/11100465-M.jpg', 1933, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000080', '9780804139298', 'Zero to One', 'Peter A. Thiel, Blake Masters', NULL,
	'https://covers.openlibrary.org/b/id/9002334-M.jpg', 2001, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000081', '0811204812', '人間失格', '太宰 治', NULL,
	'https://covers.openlibrary.org/b/id/13190147-M.jpg', 1948, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000082', '9780768413007', 'The art of persuasion', 'Bob Burg', NULL,
	'https://covers.openlibrary.org/b/id/12719675-M.jpg', 2011, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000083', '9780802473158', 'The Five Love Languages', 'Gary D. Chapman', NULL,
	'https://covers.openlibrary.org/b/id/12602983-M.jpg', 1992, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000084', '9781878424938', 'The Four Agreements', 'Don Miguel Ruiz', NULL,
	'https://covers.openlibrary.org/b/id/924521-M.jpg', 1997, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000085', '9781408855942', 'Harry Potter and the Half-Blood Prince', 'J. K. Rowling', NULL,
	'https://covers.openlibrary.org/b/id/10716273-M.jpg', 2005, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000086', '0385247125', 'Moonwalk', 'Michael Jackson', NULL,
	'https://covers.openlibrary.org/b/id/6294012-M.jpg', 1988, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000087', '9780340797662', 'The Shining', 'Stephen King', NULL,
	'https://covers.openlibrary.org/b/id/12376585-M.jpg', 1977, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000088', '1982150920', 'Cadáver exquisito', 'Agustina Bazterrica', NULL,
	'https://covers.openlibrary.org/b/id/8169547-M.jpg', 2015, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000089', '043935806X', 'Harry Potter and the Order of the Phoenix', 'J. K. Rowling', NULL,
	'https://covers.openlibrary.org/b/id/15158666-M.jpg', 2003, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000090', '0205361374', 'The world of psychology', 'Samuel E. Wood, Ellen Green Wood, Denise Boyd, Wood, Ellen R. Green Wood, WOOD WOOD, Ellen Wood', NULL,
	'https://covers.openlibrary.org/b/id/15170540-M.jpg', 1992, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000091', '0385516487', '''Salem’s Lot', 'Stephen King', NULL,
	'https://covers.openlibrary.org/b/id/14654118-M.jpg', 1975, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000092', '9780099579939', 'Fifty Shades of Grey', 'E. L. James', NULL,
	'https://covers.openlibrary.org/b/id/12648183-M.jpg', 2000, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000093', '9781481456203', 'The Summer I Turned Pretty Trilogy', 'Jenny Han', NULL,
	'https://covers.openlibrary.org/b/id/7719210-M.jpg', 2009, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000094', '9781577314806', 'The Power of Now', 'Eckhart Tolle', NULL,
	'https://covers.openlibrary.org/b/id/551262-M.jpg', 1997, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000095', '0758796005', 'The Perks of Being a Wallflower', 'Stephen Chbosky', NULL,
	'https://covers.openlibrary.org/b/id/14315052-M.jpg', 1999, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000096', '9780753189603', 'The Song of Achilles', 'Madeline Miller', NULL,
	'https://covers.openlibrary.org/b/id/7098465-M.jpg', 2011, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000097', '9780536631237', 'Biology', 'Neil Alexander Campbell', NULL,
	'https://covers.openlibrary.org/b/id/581911-M.jpg', 1987, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000098', '074754624X', 'Harry Potter and the Goblet of Fire', 'J. K. Rowling', NULL,
	'https://covers.openlibrary.org/b/id/12059372-M.jpg', 2000, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000099', '9780553113402', 'The Exorcist', 'William Peter Blatty', NULL,
	'https://covers.openlibrary.org/b/id/12715730-M.jpg', 1971, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
),
(
	'10000000-0000-4000-8000-000000000100', '9781608612819', 'Seeking Persephone', 'Sarah M. Eden', NULL,
	'https://covers.openlibrary.org/b/id/10864164-M.jpg', 2008, 5, 5,
	'2026-07-19 00:00:00.000000', '2026-07-19 00:00:00.000000'
);
