-- Correct category-specific library-staff copy introduced by migration 202607190007.
-- Preserve descriptions changed by a librarian after the earlier migration.
UPDATE books
SET description = CASE id
	WHEN '10000000-0000-4000-8000-000000000020' THEN 'Greg Heffley records the everyday embarrassments, schemes, and family complications of middle-school life in his illustrated diary.

Cartoon illustrations, quick scenes, and comic misunderstandings make the everyday pressures of school and family immediately recognizable.

Recommended for middle-grade readers looking for funny, highly accessible realistic fiction.'
	WHEN '10000000-0000-4000-8000-000000000022' THEN 'A pilot stranded in the desert meets a young prince whose small-planet stories explore friendship, loss, responsibility, and wonder.

Its spare, fable-like form makes room for reflection, and its themes continue to invite readers of many ages.

Recommended for readers seeking a short classic with philosophical warmth.'
	WHEN '10000000-0000-4000-8000-000000000068' THEN 'A boy with a facial difference enters school for the first time, changing the way classmates, teachers, and family members see courage.

A warm, realistic school story, it asks how kindness can change a community without pretending that empathy is always easy.

Recommended for middle-grade readers, families, and classrooms ready to talk about belonging and compassion.'
	WHEN '10000000-0000-4000-8000-000000000072' THEN 'An orphaned governess seeks independence and love while confronting class, conscience, and the secrets of Thornfield Hall.

The novel combines a fiercely independent heroine with Gothic atmosphere, moral conflict, and a memorable portrait of desire and self-respect.

Recommended for readers who enjoy enduring classics with romance, suspense, and a strong narrative voice.'
END
WHERE id IN (
	'10000000-0000-4000-8000-000000000020',
	'10000000-0000-4000-8000-000000000022',
	'10000000-0000-4000-8000-000000000068',
	'10000000-0000-4000-8000-000000000072'
) AND description = CASE id
	WHEN '10000000-0000-4000-8000-000000000020' THEN 'Greg Heffley records the everyday embarrassments, schemes, and family complications of middle-school life in his illustrated diary.

Adventure, humor, and mythic stakes keep the story approachable while friendship and courage remain at its center.

Recommended for middle-grade readers and families who enjoy fast, imaginative series fiction.'
	WHEN '10000000-0000-4000-8000-000000000022' THEN 'A pilot stranded in the desert meets a young prince whose small-planet stories explore friendship, loss, responsibility, and wonder.

Its spare, fable-like form makes room for reflection, and its themes continue to invite readers of many ages.

Recommended for readers seeking a short classic with philosophical warmth.'
	WHEN '10000000-0000-4000-8000-000000000068' THEN 'A boy with a facial difference enters school for the first time, changing the way classmates, teachers, and family members see courage.

Adventure, humor, and mythic stakes keep the story approachable while friendship and courage remain at its center.

Recommended for middle-grade readers and families who enjoy fast, imaginative series fiction.'
	WHEN '10000000-0000-4000-8000-000000000072' THEN 'An orphaned governess seeks independence and love while confronting class, conscience, and the secrets of Thornfield Hall.

Its spare, fable-like form makes room for reflection, and its themes continue to invite readers of many ages.

Recommended for readers seeking a short classic with philosophical warmth.'
ELSE description
END;
