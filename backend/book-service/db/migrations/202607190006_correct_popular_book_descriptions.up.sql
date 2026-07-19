-- Remove unverified first-publication claims from the earlier generated descriptions.
-- Preserve any librarian-authored description that differs from the generated value.
UPDATE books
SET description = CASE id
	WHEN '10000000-0000-4000-8000-000000000011' THEN 'The Hunger Games is a popular book by Suzanne Collins. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000012' THEN 'The Laws of Human Nature is a popular book by Robert Greene. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000013' THEN 'Thinking, fast and slow is a popular book by Daniel Kahneman. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000014' THEN 'Nineteen Eighty-Four is a popular book by George Orwell. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000015' THEN 'The psychology of money is a popular book by Henry Clay Lindgren. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000016' THEN 'The Subtle Art of Not Giving a F*ck is a popular book by Mark Manson. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000017' THEN 'The Lightning Thief is a popular book by Rick Riordan. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000018' THEN 'A Little Life is a popular book by Hanya Yanagihara. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000019' THEN 'Fire & Blood is a popular book by George R. R. Martin. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000020' THEN 'Diary of a Wimpy Kid is a popular book by Jeff Kinney. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000021' THEN 'Shatter Me is a popular book by Tahereh Mafi. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000022' THEN 'Le petit prince is a popular book by Antoine de Saint-Exupéry. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000023' THEN 'A Good Girl''s Guide to Murder is a popular book by Holly Jackson. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000024' THEN 'The Art of Seduction is a popular book by Robert Greene, Joost Elffers. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000025' THEN 'The Power of Your Subconscious Mind is a popular book by Joseph Murphy. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000026' THEN '嫌われる勇気 is a popular book by Ichirō Kishimi, Fumitake Koga. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000027' THEN 'Project Hail Mary is a popular book by Andy Weir. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000028' THEN 'Twilight is a popular book by Stephenie Meyer. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000029' THEN 'It Ends With Us is a popular book by Colleen Hoover. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000030' THEN 'The Mountain Is You is a popular book by Brianna Wiest. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000031' THEN 'The Mistake is a popular book by Elle Kennedy. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000032' THEN 'Harry Potter and the Chamber of Secrets is a popular book by J. K. Rowling. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000033' THEN 'Deep Work is a popular book by Cal Newport. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000034' THEN 'Sapiens is a popular book by Yuval Noah Harari. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000035' THEN 'The Richest Man in Babylon is a popular book by George S. Clason. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000036' THEN 'Psychology is a popular book by Robert A. Baron. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000037' THEN 'Ikigai is a popular book by Héctor García, Francesc Miralles. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000038' THEN 'Mindset is a popular book by Carol S. Dweck. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000039' THEN 'Harry Potter and the Prisoner of Azkaban is a popular book by J. K. Rowling. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000040' THEN 'The Hobbit is a popular book by J.R.R. Tolkien. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000041' THEN 'It is a popular book by Stephen King. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000042' THEN 'The Kite Runner is a popular book by Khaled Hosseini. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000043' THEN 'How to Talk to Anyone is a popular book by Leil Lowndes. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000044' THEN 'Trading in the zone is a popular book by Mark Douglas. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000045' THEN 'Red Queen is a popular book by Victoria Aveyard. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000046' THEN 'The Battle of the Labyrinth is a popular book by Rick Riordan. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000047' THEN 'The Last Olympian is a popular book by Rick Riordan. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000048' THEN 'Red Rising is a popular book by Pierce Brown. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000049' THEN 'Suicide Med is a popular book by Freida McFadden. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000050' THEN 'Rich Dad, Poor Dad for Teens is a popular book by Robert T. Kiyosaki, Sharon L. Lechter. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000051' THEN 'The Hired Girl is a popular book by Laura Amy Schlitz. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000052' THEN 'The Power of Logic is a popular book by Frances Howard-Snyder, Daniel Howard-Snyder, Ryan Wasserman. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000053' THEN 'The Lost Hero is a popular book by Rick Riordan. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000054' THEN 'The Fault in Our Stars is a popular book by John Green. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000055' THEN '$100M Offers is a popular book by Alex Hormozi. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000056' THEN 'Brain Damage is a popular book by Freida McFadden. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000057' THEN 'God of Malice is a popular book by Rina Kent. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000058' THEN 'Forbidden knowledge is a popular book by Roger Shattuck. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000059' THEN 'A Thousand Splendid Suns is a popular book by Khaled Hosseini. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000060' THEN 'Corrupt is a popular book by Penelope Douglas. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000061' THEN 'Influence is a popular book by Robert B. Cialdini. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000062' THEN 'Misery is a popular book by Stephen King. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000063' THEN 'The Score is a popular book by Elle Kennedy. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000064' THEN 'Girl in Pieces is a popular book by Kathleen Glasgow. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000065' THEN 'Surrounded by Idiots is a popular book by Thomas Erikson. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000066' THEN 'A Court of Thorns and Roses is a popular book by Sarah J. Maas. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000067' THEN 'The Titan''s Curse is a popular book by Rick Riordan. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000068' THEN 'Wonder is a popular book by R. J. Palacio. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000069' THEN 'Cien años de soledad is a popular book by Gabriel García Márquez. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000070' THEN 'The Body Keeps the Score is a popular book by Bessel van der Kolk. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000071' THEN 'Emotional Intelligence is a popular book by Daniel Goleman. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000072' THEN 'Jane Eyre is a popular book by Charlotte Brontë. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000073' THEN 'The 48 Laws of Power Pivotal Points -The Pivotal Guide to Robert Greene''s Celebrated Book is a popular book by Pivotal Point Papers. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000074' THEN 'The Secret is a popular book by Rhonda Byrne. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000075' THEN 'I, Robot is a popular book by Isaac Asimov. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000076' THEN 'Harry Potter and the Deathly Hallows is a popular book by J. K. Rowling. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000077' THEN 'Diary of a Wimpy Kid is a popular book by Jeff Kinney. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000078' THEN 'The Sea of Monsters is a popular book by Rick Riordan. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000079' THEN 'Murder on the Orient Express is a popular book by Agatha Christie. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000080' THEN 'Zero to One is a popular book by Peter A. Thiel, Blake Masters. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000081' THEN '人間失格 is a popular book by 太宰 治. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000082' THEN 'The art of persuasion is a popular book by Bob Burg. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000083' THEN 'The Five Love Languages is a popular book by Gary D. Chapman. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000084' THEN 'The Four Agreements is a popular book by Don Miguel Ruiz. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000085' THEN 'Harry Potter and the Half-Blood Prince is a popular book by J. K. Rowling. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000086' THEN 'Moonwalk is a popular book by Michael Jackson. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000087' THEN 'The Shining is a popular book by Stephen King. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000088' THEN 'Cadáver exquisito is a popular book by Agustina Bazterrica. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000089' THEN 'Harry Potter and the Order of the Phoenix is a popular book by J. K. Rowling. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000090' THEN 'The world of psychology is a popular book by Samuel E. Wood, Ellen Green Wood, Denise Boyd, Wood, Ellen R. Green Wood, WOOD WOOD, Ellen Wood. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000091' THEN '''Salem’s Lot is a popular book by Stephen King. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000092' THEN 'Fifty Shades of Grey is a popular book by E. L. James. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000093' THEN 'The Summer I Turned Pretty Trilogy is a popular book by Jenny Han. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000094' THEN 'The Power of Now is a popular book by Eckhart Tolle. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000095' THEN 'The Perks of Being a Wallflower is a popular book by Stephen Chbosky. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000096' THEN 'The Song of Achilles is a popular book by Madeline Miller. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000097' THEN 'Biology is a popular book by Neil Alexander Campbell. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000098' THEN 'Harry Potter and the Goblet of Fire is a popular book by J. K. Rowling. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000099' THEN 'The Exorcist is a popular book by William Peter Blatty. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000100' THEN 'Seeking Persephone is a popular book by Sarah M. Eden. This library edition is included in the Open Library monthly-trending catalog.'
END
WHERE id IN (
	'10000000-0000-4000-8000-000000000011',
	'10000000-0000-4000-8000-000000000012',
	'10000000-0000-4000-8000-000000000013',
	'10000000-0000-4000-8000-000000000014',
	'10000000-0000-4000-8000-000000000015',
	'10000000-0000-4000-8000-000000000016',
	'10000000-0000-4000-8000-000000000017',
	'10000000-0000-4000-8000-000000000018',
	'10000000-0000-4000-8000-000000000019',
	'10000000-0000-4000-8000-000000000020',
	'10000000-0000-4000-8000-000000000021',
	'10000000-0000-4000-8000-000000000022',
	'10000000-0000-4000-8000-000000000023',
	'10000000-0000-4000-8000-000000000024',
	'10000000-0000-4000-8000-000000000025',
	'10000000-0000-4000-8000-000000000026',
	'10000000-0000-4000-8000-000000000027',
	'10000000-0000-4000-8000-000000000028',
	'10000000-0000-4000-8000-000000000029',
	'10000000-0000-4000-8000-000000000030',
	'10000000-0000-4000-8000-000000000031',
	'10000000-0000-4000-8000-000000000032',
	'10000000-0000-4000-8000-000000000033',
	'10000000-0000-4000-8000-000000000034',
	'10000000-0000-4000-8000-000000000035',
	'10000000-0000-4000-8000-000000000036',
	'10000000-0000-4000-8000-000000000037',
	'10000000-0000-4000-8000-000000000038',
	'10000000-0000-4000-8000-000000000039',
	'10000000-0000-4000-8000-000000000040',
	'10000000-0000-4000-8000-000000000041',
	'10000000-0000-4000-8000-000000000042',
	'10000000-0000-4000-8000-000000000043',
	'10000000-0000-4000-8000-000000000044',
	'10000000-0000-4000-8000-000000000045',
	'10000000-0000-4000-8000-000000000046',
	'10000000-0000-4000-8000-000000000047',
	'10000000-0000-4000-8000-000000000048',
	'10000000-0000-4000-8000-000000000049',
	'10000000-0000-4000-8000-000000000050',
	'10000000-0000-4000-8000-000000000051',
	'10000000-0000-4000-8000-000000000052',
	'10000000-0000-4000-8000-000000000053',
	'10000000-0000-4000-8000-000000000054',
	'10000000-0000-4000-8000-000000000055',
	'10000000-0000-4000-8000-000000000056',
	'10000000-0000-4000-8000-000000000057',
	'10000000-0000-4000-8000-000000000058',
	'10000000-0000-4000-8000-000000000059',
	'10000000-0000-4000-8000-000000000060',
	'10000000-0000-4000-8000-000000000061',
	'10000000-0000-4000-8000-000000000062',
	'10000000-0000-4000-8000-000000000063',
	'10000000-0000-4000-8000-000000000064',
	'10000000-0000-4000-8000-000000000065',
	'10000000-0000-4000-8000-000000000066',
	'10000000-0000-4000-8000-000000000067',
	'10000000-0000-4000-8000-000000000068',
	'10000000-0000-4000-8000-000000000069',
	'10000000-0000-4000-8000-000000000070',
	'10000000-0000-4000-8000-000000000071',
	'10000000-0000-4000-8000-000000000072',
	'10000000-0000-4000-8000-000000000073',
	'10000000-0000-4000-8000-000000000074',
	'10000000-0000-4000-8000-000000000075',
	'10000000-0000-4000-8000-000000000076',
	'10000000-0000-4000-8000-000000000077',
	'10000000-0000-4000-8000-000000000078',
	'10000000-0000-4000-8000-000000000079',
	'10000000-0000-4000-8000-000000000080',
	'10000000-0000-4000-8000-000000000081',
	'10000000-0000-4000-8000-000000000082',
	'10000000-0000-4000-8000-000000000083',
	'10000000-0000-4000-8000-000000000084',
	'10000000-0000-4000-8000-000000000085',
	'10000000-0000-4000-8000-000000000086',
	'10000000-0000-4000-8000-000000000087',
	'10000000-0000-4000-8000-000000000088',
	'10000000-0000-4000-8000-000000000089',
	'10000000-0000-4000-8000-000000000090',
	'10000000-0000-4000-8000-000000000091',
	'10000000-0000-4000-8000-000000000092',
	'10000000-0000-4000-8000-000000000093',
	'10000000-0000-4000-8000-000000000094',
	'10000000-0000-4000-8000-000000000095',
	'10000000-0000-4000-8000-000000000096',
	'10000000-0000-4000-8000-000000000097',
	'10000000-0000-4000-8000-000000000098',
	'10000000-0000-4000-8000-000000000099',
	'10000000-0000-4000-8000-000000000100'
) AND description = CASE id
	WHEN '10000000-0000-4000-8000-000000000011' THEN 'First published in 2008, The Hunger Games is a popular book by Suzanne Collins. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000012' THEN 'First published in 2018, The Laws of Human Nature is a popular book by Robert Greene. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000013' THEN 'First published in 2011, Thinking, fast and slow is a popular book by Daniel Kahneman. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000014' THEN 'First published in 1949, Nineteen Eighty-Four is a popular book by George Orwell. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000015' THEN 'First published in 1991, The psychology of money is a popular book by Henry Clay Lindgren. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000016' THEN 'First published in 2016, The Subtle Art of Not Giving a F*ck is a popular book by Mark Manson. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000017' THEN 'First published in 2005, The Lightning Thief is a popular book by Rick Riordan. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000018' THEN 'First published in 2008, A Little Life is a popular book by Hanya Yanagihara. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000019' THEN 'First published in 2014, Fire & Blood is a popular book by George R. R. Martin. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000020' THEN 'First published in 2007, Diary of a Wimpy Kid is a popular book by Jeff Kinney. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000021' THEN 'First published in 2011, Shatter Me is a popular book by Tahereh Mafi. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000022' THEN 'First published in 1943, Le petit prince is a popular book by Antoine de Saint-Exupéry. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000023' THEN 'First published in 2019, A Good Girl''s Guide to Murder is a popular book by Holly Jackson. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000024' THEN 'First published in 2001, The Art of Seduction is a popular book by Robert Greene, Joost Elffers. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000025' THEN 'First published in 1963, The Power of Your Subconscious Mind is a popular book by Joseph Murphy. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000026' THEN 'First published in 2013, 嫌われる勇気 is a popular book by Ichirō Kishimi, Fumitake Koga. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000027' THEN 'First published in 2021, Project Hail Mary is a popular book by Andy Weir. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000028' THEN 'First published in 2005, Twilight is a popular book by Stephenie Meyer. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000029' THEN 'First published in 2012, It Ends With Us is a popular book by Colleen Hoover. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000030' THEN 'First published in 2020, The Mountain Is You is a popular book by Brianna Wiest. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000031' THEN 'First published in 2015, The Mistake is a popular book by Elle Kennedy. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000032' THEN 'First published in 1998, Harry Potter and the Chamber of Secrets is a popular book by J. K. Rowling. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000033' THEN 'First published in 2016, Deep Work is a popular book by Cal Newport. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000034' THEN 'First published in 2011, Sapiens is a popular book by Yuval Noah Harari. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000035' THEN 'First published in 1926, The Richest Man in Babylon is a popular book by George S. Clason. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000036' THEN 'First published in 1977, Psychology is a popular book by Robert A. Baron. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000037' THEN 'First published in 2016, Ikigai is a popular book by Héctor García, Francesc Miralles. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000038' THEN 'First published in 2006, Mindset is a popular book by Carol S. Dweck. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000039' THEN 'First published in 1999, Harry Potter and the Prisoner of Azkaban is a popular book by J. K. Rowling. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000040' THEN 'First published in 1703, The Hobbit is a popular book by J.R.R. Tolkien. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000041' THEN 'First published in 1986, It is a popular book by Stephen King. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000042' THEN 'First published in 2003, The Kite Runner is a popular book by Khaled Hosseini. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000043' THEN 'First published in 1999, How to Talk to Anyone is a popular book by Leil Lowndes. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000044' THEN 'First published in 2001, Trading in the zone is a popular book by Mark Douglas. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000045' THEN 'First published in 2015, Red Queen is a popular book by Victoria Aveyard. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000046' THEN 'First published in 2005, The Battle of the Labyrinth is a popular book by Rick Riordan. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000047' THEN 'First published in 2009, The Last Olympian is a popular book by Rick Riordan. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000048' THEN 'First published in 2014, Red Rising is a popular book by Pierce Brown. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000049' THEN 'First published in 2014, Suicide Med is a popular book by Freida McFadden. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000050' THEN 'First published in 2004, Rich Dad, Poor Dad for Teens is a popular book by Robert T. Kiyosaki, Sharon L. Lechter. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000051' THEN 'First published in 2015, The Hired Girl is a popular book by Laura Amy Schlitz. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000052' THEN 'First published in 2009, The Power of Logic is a popular book by Frances Howard-Snyder, Daniel Howard-Snyder, Ryan Wasserman. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000053' THEN 'First published in 2010, The Lost Hero is a popular book by Rick Riordan. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000054' THEN 'First published in 2010, The Fault in Our Stars is a popular book by John Green. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000055' THEN 'First published in 2021, $100M Offers is a popular book by Alex Hormozi. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000056' THEN 'First published in 2016, Brain Damage is a popular book by Freida McFadden. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000057' THEN 'First published in 2022, God of Malice is a popular book by Rina Kent. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000058' THEN 'First published in 1996, Forbidden knowledge is a popular book by Roger Shattuck. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000059' THEN 'First published in 2007, A Thousand Splendid Suns is a popular book by Khaled Hosseini. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000060' THEN 'First published in 2015, Corrupt is a popular book by Penelope Douglas. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000061' THEN 'First published in 1983, Influence is a popular book by Robert B. Cialdini. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000062' THEN 'First published in 1978, Misery is a popular book by Stephen King. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000063' THEN 'First published in 2016, The Score is a popular book by Elle Kennedy. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000064' THEN 'First published in 2000, Girl in Pieces is a popular book by Kathleen Glasgow. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000065' THEN 'First published in 2014, Surrounded by Idiots is a popular book by Thomas Erikson. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000066' THEN 'First published in 2013, A Court of Thorns and Roses is a popular book by Sarah J. Maas. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000067' THEN 'First published in 2007, The Titan''s Curse is a popular book by Rick Riordan. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000068' THEN 'First published in 2001, Wonder is a popular book by R. J. Palacio. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000069' THEN 'First published in 1967, Cien años de soledad is a popular book by Gabriel García Márquez. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000070' THEN 'First published in 2014, The Body Keeps the Score is a popular book by Bessel van der Kolk. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000071' THEN 'First published in 1995, Emotional Intelligence is a popular book by Daniel Goleman. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000072' THEN 'First published in 1847, Jane Eyre is a popular book by Charlotte Brontë. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000073' THEN 'First published in 2014, The 48 Laws of Power Pivotal Points -The Pivotal Guide to Robert Greene''s Celebrated Book is a popular book by Pivotal Point Papers. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000074' THEN 'First published in 2000, The Secret is a popular book by Rhonda Byrne. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000075' THEN 'First published in 1950, I, Robot is a popular book by Isaac Asimov. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000076' THEN 'First published in 2007, Harry Potter and the Deathly Hallows is a popular book by J. K. Rowling. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000077' THEN 'First published in 2017, Diary of a Wimpy Kid is a popular book by Jeff Kinney. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000078' THEN 'First published in 2005, The Sea of Monsters is a popular book by Rick Riordan. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000079' THEN 'First published in 1933, Murder on the Orient Express is a popular book by Agatha Christie. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000080' THEN 'First published in 2001, Zero to One is a popular book by Peter A. Thiel, Blake Masters. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000081' THEN 'First published in 1948, 人間失格 is a popular book by 太宰 治. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000082' THEN 'First published in 2011, The art of persuasion is a popular book by Bob Burg. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000083' THEN 'First published in 1992, The Five Love Languages is a popular book by Gary D. Chapman. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000084' THEN 'First published in 1997, The Four Agreements is a popular book by Don Miguel Ruiz. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000085' THEN 'First published in 2005, Harry Potter and the Half-Blood Prince is a popular book by J. K. Rowling. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000086' THEN 'First published in 1988, Moonwalk is a popular book by Michael Jackson. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000087' THEN 'First published in 1977, The Shining is a popular book by Stephen King. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000088' THEN 'First published in 2015, Cadáver exquisito is a popular book by Agustina Bazterrica. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000089' THEN 'First published in 2003, Harry Potter and the Order of the Phoenix is a popular book by J. K. Rowling. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000090' THEN 'First published in 1992, The world of psychology is a popular book by Samuel E. Wood, Ellen Green Wood, Denise Boyd, Wood, Ellen R. Green Wood, WOOD WOOD, Ellen Wood. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000091' THEN 'First published in 1975, ''Salem’s Lot is a popular book by Stephen King. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000092' THEN 'First published in 2000, Fifty Shades of Grey is a popular book by E. L. James. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000093' THEN 'First published in 2009, The Summer I Turned Pretty Trilogy is a popular book by Jenny Han. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000094' THEN 'First published in 1997, The Power of Now is a popular book by Eckhart Tolle. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000095' THEN 'First published in 1999, The Perks of Being a Wallflower is a popular book by Stephen Chbosky. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000096' THEN 'First published in 2011, The Song of Achilles is a popular book by Madeline Miller. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000097' THEN 'First published in 1987, Biology is a popular book by Neil Alexander Campbell. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000098' THEN 'First published in 2000, Harry Potter and the Goblet of Fire is a popular book by J. K. Rowling. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000099' THEN 'First published in 1971, The Exorcist is a popular book by William Peter Blatty. This library edition is included in the Open Library monthly-trending catalog.'
	WHEN '10000000-0000-4000-8000-000000000100' THEN 'First published in 2008, Seeking Persephone is a popular book by Sarah M. Eden. This library edition is included in the Open Library monthly-trending catalog.'
ELSE description
END;
