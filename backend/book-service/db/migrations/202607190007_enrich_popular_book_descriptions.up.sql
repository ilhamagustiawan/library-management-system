-- Replace generic seed text with original multi-paragraph library-staff descriptions.
-- Preserve any librarian-authored description that differs from the prior generated text.
UPDATE books
SET description = CASE id
	WHEN '10000000-0000-4000-8000-000000000011' THEN 'Katniss Everdeen is forced to represent her impoverished district in a televised contest where survival is the only prize.

The novel pairs rapid pacing with questions about power, spectacle, inequality, and moral compromise.

Recommended for readers who enjoy political speculative fiction and survival stories.'
	WHEN '10000000-0000-4000-8000-000000000012' THEN 'Robert Greene examines recurring patterns in motivation, status, empathy, and conflict to help readers notice the forces shaping social behavior.

Its ideas are best read as prompts for reflection and discussion rather than as a substitute for professional advice.

Recommended for readers interested in behavior, communication, and practical self-understanding.'
	WHEN '10000000-0000-4000-8000-000000000013' THEN 'Daniel Kahneman introduces the fast, intuitive mind and the slower, deliberate mind, showing how each affects judgment and choice.

Its ideas are best read as prompts for reflection and discussion rather than as a substitute for professional advice.

Recommended for readers interested in behavior, communication, and practical self-understanding.'
	WHEN '10000000-0000-4000-8000-000000000014' THEN 'In a surveillance state where language and memory are controlled, Winston Smith begins a dangerous private rebellion against an all-powerful regime.

The novel pairs rapid pacing with questions about power, spectacle, inequality, and moral compromise.

Recommended for readers who enjoy political speculative fiction and survival stories.'
	WHEN '10000000-0000-4000-8000-000000000015' THEN 'This psychology-focused volume considers how attitudes, motives, and behavior shape the way people think about money.

Its ideas are best read as prompts for reflection and discussion rather than as a substitute for professional advice.

Recommended for readers interested in behavior, communication, and practical self-understanding.'
	WHEN '10000000-0000-4000-8000-000000000016' THEN 'Mark Manson argues for choosing values carefully, accepting limits, and directing energy toward the problems that genuinely matter.

The approach is reflective and practical, offering concepts readers can test against their own routines and priorities.

Recommended for readers looking for a thoughtful starting point for personal change.'
	WHEN '10000000-0000-4000-8000-000000000017' THEN 'Percy Jackson discovers he is tied to the Greek gods and is sent on a fast-moving quest to prevent a war among them.

Adventure, humor, and mythic stakes keep the story approachable while friendship and courage remain at its center.

Recommended for middle-grade readers and families who enjoy fast, imaginative series fiction.'
	WHEN '10000000-0000-4000-8000-000000000018' THEN 'Four college friends remain connected across decades as one of them struggles with the lasting effects of trauma and self-destruction.

The book gives close attention to character, memory, and emotional consequence, and it may be a demanding read for some audiences.

Recommended for readers who appreciate character-driven fiction and discussion-rich themes.'
	WHEN '10000000-0000-4000-8000-000000000019' THEN 'George R. R. Martin presents a history of House Targaryen, its dragons, rivalries, and civil war long before the events of A Game of Thrones.

Its world-building, family history, and competing claims to power reward readers who enjoy lore as much as action.

Recommended for fantasy readers drawn to epic histories and political conflict.'
	WHEN '10000000-0000-4000-8000-000000000020' THEN 'Greg Heffley records the everyday embarrassments, schemes, and family complications of middle-school life in his illustrated diary.

Cartoon illustrations, quick scenes, and comic misunderstandings make the everyday pressures of school and family immediately recognizable.

Recommended for middle-grade readers looking for funny, highly accessible realistic fiction.'
	WHEN '10000000-0000-4000-8000-000000000021' THEN 'Juliette possesses a lethal touch and must decide whether to remain a weapon for a controlling regime or claim her own future.

The story combines a coming-of-age voice with rebellion, danger, and a high-emotion romantic thread.

Recommended for teen and adult readers who enjoy dystopian fantasy with a strong central heroine.'
	WHEN '10000000-0000-4000-8000-000000000022' THEN 'A pilot stranded in the desert meets a young prince whose small-planet stories explore friendship, loss, responsibility, and wonder.

Its spare, fable-like form makes room for reflection, and its themes continue to invite readers of many ages.

Recommended for readers seeking a short classic with philosophical warmth.'
	WHEN '10000000-0000-4000-8000-000000000023' THEN 'A student revisits a closed murder case for a school project and discovers that someone is still determined to keep the truth buried.

Clues, shifting suspicions, and escalating stakes make it a page-turner while keeping its focus on a young investigator.

Recommended for readers who like contemporary YA mysteries and twist-driven plots.'
	WHEN '10000000-0000-4000-8000-000000000024' THEN 'Robert Greene surveys historical figures and social strategies to examine attraction, charisma, power, and manipulation.

Its ideas are best read as prompts for reflection and discussion rather than as a substitute for professional advice.

Recommended for readers interested in behavior, communication, and practical self-understanding.'
	WHEN '10000000-0000-4000-8000-000000000025' THEN 'Joseph Murphy presents a spiritual self-help approach that links belief, habit, and visualization to personal change.

The approach is reflective and practical, offering concepts readers can test against their own routines and priorities.

Recommended for readers looking for a thoughtful starting point for personal change.'
	WHEN '10000000-0000-4000-8000-000000000026' THEN 'Using ideas associated with Adlerian psychology, this dialogue explores freedom from approval-seeking and the responsibility of choosing one’s path.

Its ideas are best read as prompts for reflection and discussion rather than as a substitute for professional advice.

Recommended for readers interested in behavior, communication, and practical self-understanding.'
	WHEN '10000000-0000-4000-8000-000000000027' THEN 'A lone astronaut awakens far from Earth with a mission to solve a threat to humanity and an unexpected ally at his side.

Scientific problem-solving and escalating stakes drive the narrative, balancing technical ideas with an accessible sense of adventure.

Recommended for readers who enjoy hopeful, high-concept science fiction.'
	WHEN '10000000-0000-4000-8000-000000000028' THEN 'Bella Swan’s move to a rainy small town leads to a dangerous romance with Edward Cullen, a vampire struggling to resist his instincts.

At its center is an emotionally charged relationship shaped by identity, loyalty, and the risks of first love.

Recommended for readers drawn to romantic YA fiction with strong atmosphere.'
	WHEN '10000000-0000-4000-8000-000000000029' THEN 'Lily Bloom begins a new relationship while confronting the patterns of violence and loyalty that have shaped her family history.

The relationship develops alongside questions of trust, consent, and the choices people make after being hurt.

Recommended for adult readers who enjoy contemporary romance with emotional conflict.'
	WHEN '10000000-0000-4000-8000-000000000030' THEN 'Brianna Wiest reflects on self-sabotage, emotional patterns, and the difficult work of becoming an active participant in one’s own life.

The approach is reflective and practical, offering concepts readers can test against their own routines and priorities.

Recommended for readers looking for a thoughtful starting point for personal change.'
	WHEN '10000000-0000-4000-8000-000000000031' THEN 'A college hockey player’s past mistake and a new relationship complicate his effort to rebuild trust and move forward.

The relationship develops alongside questions of trust, consent, and the choices people make after being hurt.

Recommended for adult readers who enjoy contemporary romance with emotional conflict.'
	WHEN '10000000-0000-4000-8000-000000000032' THEN 'Harry returns to Hogwarts as a hidden chamber is opened and students are attacked, forcing him to investigate the school’s past.

Adventure, humor, and mythic stakes keep the story approachable while friendship and courage remain at its center.

Recommended for middle-grade readers and families who enjoy fast, imaginative series fiction.'
	WHEN '10000000-0000-4000-8000-000000000033' THEN 'Cal Newport makes the case for sustained, distraction-free concentration as a valuable skill for learning, creating, and professional work.

The book challenges the constant-fragmentation model of work and encourages readers to examine the conditions that support focus.

Recommended for readers looking for a disciplined approach to study or creative work.'
	WHEN '10000000-0000-4000-8000-000000000034' THEN 'Yuval Noah Harari traces the large-scale stories, institutions, and shared beliefs that shaped human societies from prehistory to the present.

It invites readers to compare competing ideas and historical evidence, making it well suited to conversation and critical reading.

Recommended for readers interested in big ideas, culture, and public debate.'
	WHEN '10000000-0000-4000-8000-000000000035' THEN 'A collection of parables set in ancient Babylon offers memorable lessons about saving, debt, opportunity, and long-term financial discipline.

Its examples emphasize habits and decision-making, encouraging readers to adapt the ideas to their own circumstances.

Recommended for readers seeking approachable personal-finance or trading perspectives.'
	WHEN '10000000-0000-4000-8000-000000000036' THEN 'This introductory textbook surveys major concepts, research methods, and debates in the study of human thought and behavior.

Definitions, examples, and structured explanations make it useful for study, review, and classroom discussion.

Recommended for students and general readers building a foundation in the subject.'
	WHEN '10000000-0000-4000-8000-000000000037' THEN 'Drawing on Japanese ideas about purpose and daily life, the authors consider small routines, community, and meaning in later years.

The approach is reflective and practical, offering concepts readers can test against their own routines and priorities.

Recommended for readers looking for a thoughtful starting point for personal change.'
	WHEN '10000000-0000-4000-8000-000000000038' THEN 'Carol Dweck explains how beliefs about ability can influence persistence, learning, feedback, and response to setbacks.

Its ideas are best read as prompts for reflection and discussion rather than as a substitute for professional advice.

Recommended for readers interested in behavior, communication, and practical self-understanding.'
	WHEN '10000000-0000-4000-8000-000000000039' THEN 'A mysterious escaped prisoner, time travel, and hidden truths about Harry’s parents reshape Harry’s third year at Hogwarts.

Adventure, humor, and mythic stakes keep the story approachable while friendship and courage remain at its center.

Recommended for middle-grade readers and families who enjoy fast, imaginative series fiction.'
	WHEN '10000000-0000-4000-8000-000000000040' THEN 'Bilbo Baggins leaves his comfortable home for a dragon-guarded treasure and a journey that introduces him to courage and adventure.

Its world-building, family history, and competing claims to power reward readers who enjoy lore as much as action.

Recommended for fantasy readers drawn to epic histories and political conflict.'
	WHEN '10000000-0000-4000-8000-000000000041' THEN 'A group of friends confronts an ancient evil that takes the form of their deepest fears in the town of Derry, Maine.

Tension builds through atmosphere, vulnerability, and the fear that familiar places or people may no longer be safe.

Recommended for readers who enjoy unsettling suspense and psychological horror.'
	WHEN '10000000-0000-4000-8000-000000000042' THEN 'Amir revisits a childhood betrayal and returns to Afghanistan seeking a chance to make amends amid political upheaval.

The book gives close attention to character, memory, and emotional consequence, and it may be a demanding read for some audiences.

Recommended for readers who appreciate character-driven fiction and discussion-rich themes.'
	WHEN '10000000-0000-4000-8000-000000000043' THEN 'Leil Lowndes offers practical conversation techniques for making introductions, listening closely, and communicating with confidence.

The approach is reflective and practical, offering concepts readers can test against their own routines and priorities.

Recommended for readers looking for a thoughtful starting point for personal change.'
	WHEN '10000000-0000-4000-8000-000000000044' THEN 'Mark Douglas examines the mental discipline, risk awareness, and consistency required to make decisions in uncertain markets.

Its examples emphasize habits and decision-making, encouraging readers to adapt the ideas to their own circumstances.

Recommended for readers seeking approachable personal-finance or trading perspectives.'
	WHEN '10000000-0000-4000-8000-000000000045' THEN 'In a kingdom divided by blood and power, Mare Barrow discovers an ability that makes her a threat to the ruling elite.

The story combines a coming-of-age voice with rebellion, danger, and a high-emotion romantic thread.

Recommended for teen and adult readers who enjoy dystopian fantasy with a strong central heroine.'
	WHEN '10000000-0000-4000-8000-000000000046' THEN 'Percy Jackson enters a dangerous maze beneath the modern world to stop a growing threat from reaching Camp Half-Blood.

Adventure, humor, and mythic stakes keep the story approachable while friendship and courage remain at its center.

Recommended for middle-grade readers and families who enjoy fast, imaginative series fiction.'
	WHEN '10000000-0000-4000-8000-000000000047' THEN 'As the Titan lord Kronos advances, Percy and his friends prepare for a final defense of Mount Olympus and the world they know.

Adventure, humor, and mythic stakes keep the story approachable while friendship and courage remain at its center.

Recommended for middle-grade readers and families who enjoy fast, imaginative series fiction.'
	WHEN '10000000-0000-4000-8000-000000000048' THEN 'Darrow infiltrates the elite class that rules a colonized solar system and finds revolution more complicated than he imagined.

Scientific problem-solving and escalating stakes drive the narrative, balancing technical ideas with an accessible sense of adventure.

Recommended for readers who enjoy hopeful, high-concept science fiction.'
	WHEN '10000000-0000-4000-8000-000000000049' THEN 'This suspense novel centers on psychological pressure, hidden motives, and the unsettling consequences of a dangerous situation.

Short scenes, uncertainty, and reversals create a quick-reading atmosphere built around suspicion and escalating risk.

Recommended for readers who want a tense, twist-forward psychological thriller.'
	WHEN '10000000-0000-4000-8000-000000000050' THEN 'A teen-focused adaptation of the Rich Dad philosophy introduces money habits, assets, goals, and choices in an accessible voice.

Its examples emphasize habits and decision-making, encouraging readers to adapt the ideas to their own circumstances.

Recommended for readers seeking approachable personal-finance or trading perspectives.'
	WHEN '10000000-0000-4000-8000-000000000051' THEN 'A young woman leaves an unhappy home for a position in Baltimore, where work, education, faith, and independence reshape her plans.

Period detail and social expectations frame the characters’ choices, adding context to the central emotional story.

Recommended for readers who enjoy historical settings, relationships, and moral dilemmas.'
	WHEN '10000000-0000-4000-8000-000000000052' THEN 'This textbook develops tools for evaluating arguments, recognizing fallacies, and constructing clear, defensible reasoning.

Definitions, examples, and structured explanations make it useful for study, review, and classroom discussion.

Recommended for students and general readers building a foundation in the subject.'
	WHEN '10000000-0000-4000-8000-000000000053' THEN 'New demigod heroes join a quest that expands the world of Percy Jackson and tests old alliances against a new prophecy.

Adventure, humor, and mythic stakes keep the story approachable while friendship and courage remain at its center.

Recommended for middle-grade readers and families who enjoy fast, imaginative series fiction.'
	WHEN '10000000-0000-4000-8000-000000000054' THEN 'Two teenagers living with cancer meet in a support group and build a relationship that is funny, tender, and painfully finite.

The story treats adolescence as a time of intense feeling, change, and connection, with themes that may be difficult for some readers.

Recommended for teen and adult readers who value emotionally candid coming-of-age fiction.'
	WHEN '10000000-0000-4000-8000-000000000055' THEN 'Alex Hormozi explains a business framework for designing offers that make a product’s value clear and compelling to customers.

The author emphasizes frameworks and choices that readers can question, adapt, and test in their own work.

Recommended for entrepreneurs and readers interested in strategy, sales, or innovation.'
	WHEN '10000000-0000-4000-8000-000000000056' THEN 'This suspense novel follows a character facing unsettling clues, fractured trust, and a situation that grows more dangerous by the page.

Short scenes, uncertainty, and reversals create a quick-reading atmosphere built around suspicion and escalating risk.

Recommended for readers who want a tense, twist-forward psychological thriller.'
	WHEN '10000000-0000-4000-8000-000000000057' THEN 'A dark college romance unfolds around a secretive elite circle, rivalries, and a relationship shaped by power and obsession.

The book explores intense attraction and unequal power in a deliberately dark register; readers may wish to check content notes.

Recommended for adult romance readers comfortable with darker themes and morally complicated characters.'
	WHEN '10000000-0000-4000-8000-000000000058' THEN 'Roger Shattuck considers the cultural cost of refusing difficult ideas and the debates surrounding taboo, censorship, and inquiry.

It invites readers to compare competing ideas and historical evidence, making it well suited to conversation and critical reading.

Recommended for readers interested in big ideas, culture, and public debate.'
	WHEN '10000000-0000-4000-8000-000000000059' THEN 'Two Afghan women form an unlikely bond while war, family expectations, and political violence transform their lives.

The book gives close attention to character, memory, and emotional consequence, and it may be a demanding read for some audiences.

Recommended for readers who appreciate character-driven fiction and discussion-rich themes.'
	WHEN '10000000-0000-4000-8000-000000000060' THEN 'A dark romance follows a young woman drawn into a world of old grudges, intimidation, and attraction with dangerous consequences.

The book explores intense attraction and unequal power in a deliberately dark register; readers may wish to check content notes.

Recommended for adult romance readers comfortable with darker themes and morally complicated characters.'
	WHEN '10000000-0000-4000-8000-000000000061' THEN 'Robert Cialdini explains the principles that make persuasion effective, from reciprocity and social proof to authority and scarcity.

Its ideas are best read as prompts for reflection and discussion rather than as a substitute for professional advice.

Recommended for readers interested in behavior, communication, and practical self-understanding.'
	WHEN '10000000-0000-4000-8000-000000000062' THEN 'A celebrated novelist is held captive by an obsessive reader who demands that he revive her favorite fictional character.

Tension builds through atmosphere, vulnerability, and the fear that familiar places or people may no longer be safe.

Recommended for readers who enjoy unsettling suspense and psychological horror.'
	WHEN '10000000-0000-4000-8000-000000000063' THEN 'A hockey star and a driven student navigate ambition, friendship, and a relationship complicated by very different plans.

The relationship develops alongside questions of trust, consent, and the choices people make after being hurt.

Recommended for adult readers who enjoy contemporary romance with emotional conflict.'
	WHEN '10000000-0000-4000-8000-000000000064' THEN 'A teenager coping with self-harm, addiction, and grief works toward recovery while learning to trust people who want to help.

The story treats adolescence as a time of intense feeling, change, and connection, with themes that may be difficult for some readers.

Recommended for teen and adult readers who value emotionally candid coming-of-age fiction.'
	WHEN '10000000-0000-4000-8000-000000000065' THEN 'Thomas Erikson presents a color-coded model of communication styles and uses it to discuss conflict, teamwork, and misunderstanding.

Its ideas are best read as prompts for reflection and discussion rather than as a substitute for professional advice.

Recommended for readers interested in behavior, communication, and practical self-understanding.'
	WHEN '10000000-0000-4000-8000-000000000066' THEN 'A human huntress is drawn into the faerie realm after killing a wolf and becomes entangled in an ancient curse and shifting loyalties.

Romance and danger move together as the heroine navigates a magical world with shifting alliances and hidden rules.

Recommended for readers who enjoy faerie fantasy, enemies-to-lovers tension, and immersive world-building.'
	WHEN '10000000-0000-4000-8000-000000000067' THEN 'Percy and his friends face a prophecy, a missing goddess, and a journey that tests their loyalty to one another.

Adventure, humor, and mythic stakes keep the story approachable while friendship and courage remain at its center.

Recommended for middle-grade readers and families who enjoy fast, imaginative series fiction.'
	WHEN '10000000-0000-4000-8000-000000000068' THEN 'A boy with a facial difference enters school for the first time, changing the way classmates, teachers, and family members see courage.

A warm, realistic school story, it asks how kindness can change a community without pretending that empathy is always easy.

Recommended for middle-grade readers, families, and classrooms ready to talk about belonging and compassion.'
	WHEN '10000000-0000-4000-8000-000000000069' THEN 'The Buendía family’s generations in the town of Macondo blend love, war, memory, and the extraordinary in a landmark saga.

The book gives close attention to character, memory, and emotional consequence, and it may be a demanding read for some audiences.

Recommended for readers who appreciate character-driven fiction and discussion-rich themes.'
	WHEN '10000000-0000-4000-8000-000000000070' THEN 'Bessel van der Kolk explains how trauma can affect body, memory, emotion, and relationships, while discussing paths toward recovery.

Its ideas are best read as prompts for reflection and discussion rather than as a substitute for professional advice.

Recommended for readers interested in behavior, communication, and practical self-understanding.'
	WHEN '10000000-0000-4000-8000-000000000071' THEN 'Daniel Goleman explores how self-awareness, empathy, motivation, and social skills contribute to decision-making and relationships.

Its ideas are best read as prompts for reflection and discussion rather than as a substitute for professional advice.

Recommended for readers interested in behavior, communication, and practical self-understanding.'
	WHEN '10000000-0000-4000-8000-000000000072' THEN 'An orphaned governess seeks independence and love while confronting class, conscience, and the secrets of Thornfield Hall.

The novel combines a fiercely independent heroine with Gothic atmosphere, moral conflict, and a memorable portrait of desire and self-respect.

Recommended for readers who enjoy enduring classics with romance, suspense, and a strong narrative voice.'
	WHEN '10000000-0000-4000-8000-000000000073' THEN 'This companion guide distills key ideas and discussion points from Robert Greene’s The 48 Laws of Power.

Its concise format is designed to support recall, discussion, and a return to the source work rather than replace it.

Recommended for readers looking for a companion or discussion aid.'
	WHEN '10000000-0000-4000-8000-000000000074' THEN 'Rhonda Byrne presents a popular self-help philosophy centered on gratitude, visualization, and the belief that thoughts shape outcomes.

The approach is reflective and practical, offering concepts readers can test against their own routines and priorities.

Recommended for readers looking for a thoughtful starting point for personal change.'
	WHEN '10000000-0000-4000-8000-000000000075' THEN 'Isaac Asimov’s linked robot stories explore artificial intelligence, ethics, and the unexpected consequences of the Three Laws of Robotics.

Scientific problem-solving and escalating stakes drive the narrative, balancing technical ideas with an accessible sense of adventure.

Recommended for readers who enjoy hopeful, high-concept science fiction.'
	WHEN '10000000-0000-4000-8000-000000000076' THEN 'Harry, Ron, and Hermione leave school to search for Horcruxes as the wizarding world falls under Voldemort’s control.

Adventure, humor, and mythic stakes keep the story approachable while friendship and courage remain at its center.

Recommended for middle-grade readers and families who enjoy fast, imaginative series fiction.'
	WHEN '10000000-0000-4000-8000-000000000077' THEN 'Greg Heffley returns with another illustrated account of school, family, friendships, and the small disasters of growing up.

Cartoon illustrations, quick scenes, and comic misunderstandings make the everyday pressures of school and family immediately recognizable.

Recommended for middle-grade readers looking for funny, highly accessible realistic fiction.'
	WHEN '10000000-0000-4000-8000-000000000078' THEN 'Percy, Annabeth, and Tyson cross dangerous waters in search of the Golden Fleece and a way to save their camp.

Adventure, humor, and mythic stakes keep the story approachable while friendship and courage remain at its center.

Recommended for middle-grade readers and families who enjoy fast, imaginative series fiction.'
	WHEN '10000000-0000-4000-8000-000000000079' THEN 'Detective Hercule Poirot investigates a murder aboard a snowbound luxury train where every passenger has a possible motive.

The contained setting and precise clue work create a classic puzzle in which every conversation matters.

Recommended for readers who enjoy traditional detective fiction and elegant locked-room puzzles.'
	WHEN '10000000-0000-4000-8000-000000000080' THEN 'Peter Thiel and Blake Masters argue that lasting businesses are built by creating something genuinely new rather than merely competing.

The author emphasizes frameworks and choices that readers can question, adapt, and test in their own work.

Recommended for entrepreneurs and readers interested in strategy, sales, or innovation.'
	WHEN '10000000-0000-4000-8000-000000000081' THEN 'Osamu Dazai’s confessional novel follows an alienated young man whose efforts to belong deepen his sense of isolation.

The book gives close attention to character, memory, and emotional consequence, and it may be a demanding read for some audiences.

Recommended for readers who appreciate character-driven fiction and discussion-rich themes.'
	WHEN '10000000-0000-4000-8000-000000000082' THEN 'Bob Burg presents practical principles for earning trust, clarifying value, and communicating persuasively without sacrificing integrity.

The author emphasizes frameworks and choices that readers can question, adapt, and test in their own work.

Recommended for entrepreneurs and readers interested in strategy, sales, or innovation.'
	WHEN '10000000-0000-4000-8000-000000000083' THEN 'Gary Chapman describes five common ways people express and receive affection, inviting readers to reflect on needs within relationships.

Its framework offers shared language for discussing care and expectation, while leaving room for individual differences.

Recommended for couples, families, and readers interested in relationship communication.'
	WHEN '10000000-0000-4000-8000-000000000084' THEN 'Don Miguel Ruiz offers a concise set of personal commitments focused on language, assumptions, integrity, and self-acceptance.

The approach is reflective and practical, offering concepts readers can test against their own routines and priorities.

Recommended for readers looking for a thoughtful starting point for personal change.'
	WHEN '10000000-0000-4000-8000-000000000085' THEN 'As the wizarding war approaches, Dumbledore shares memories that reveal Voldemort’s past and the choices Harry must face.

Adventure, humor, and mythic stakes keep the story approachable while friendship and courage remain at its center.

Recommended for middle-grade readers and families who enjoy fast, imaginative series fiction.'
	WHEN '10000000-0000-4000-8000-000000000086' THEN 'Michael Jackson recounts his musical development, creative process, and life in the public eye in this autobiographical work.

The memoir offers a first-person perspective on artistic ambition, performance, and the pressures of extraordinary fame.

Recommended for music readers and fans interested in an artist’s own account of his life.'
	WHEN '10000000-0000-4000-8000-000000000087' THEN 'A recovering writer accepts a winter job at an isolated hotel, where the building’s history and his own vulnerabilities turn frightening.

Tension builds through atmosphere, vulnerability, and the fear that familiar places or people may no longer be safe.

Recommended for readers who enjoy unsettling suspense and psychological horror.'
	WHEN '10000000-0000-4000-8000-000000000088' THEN 'In a chilling dystopia, a slaughterhouse worker begins to question a society that has normalized the consumption of human bodies.

Tension builds through atmosphere, vulnerability, and the fear that familiar places or people may no longer be safe.

Recommended for readers who enjoy unsettling suspense and psychological horror.'
	WHEN '10000000-0000-4000-8000-000000000089' THEN 'Harry forms a secret student group to prepare for a growing threat while official denial and personal loss darken his fifth year.

Adventure, humor, and mythic stakes keep the story approachable while friendship and courage remain at its center.

Recommended for middle-grade readers and families who enjoy fast, imaginative series fiction.'
	WHEN '10000000-0000-4000-8000-000000000090' THEN 'This introductory psychology text presents core theories, research, and real-world applications for students beginning the subject.

Definitions, examples, and structured explanations make it useful for study, review, and classroom discussion.

Recommended for students and general readers building a foundation in the subject.'
	WHEN '10000000-0000-4000-8000-000000000091' THEN 'A writer returns to his childhood town and finds that a vampire’s arrival is turning familiar neighbors into a hidden nightmare.

Tension builds through atmosphere, vulnerability, and the fear that familiar places or people may no longer be safe.

Recommended for readers who enjoy unsettling suspense and psychological horror.'
	WHEN '10000000-0000-4000-8000-000000000092' THEN 'An inexperienced college graduate enters an intense relationship with a wealthy entrepreneur whose desires and control challenge her boundaries.

The relationship develops alongside questions of trust, consent, and the choices people make after being hurt.

Recommended for adult readers who enjoy contemporary romance with emotional conflict.'
	WHEN '10000000-0000-4000-8000-000000000093' THEN 'Three linked summer romances follow Belly as family traditions, first love, and changing friendships complicate a beloved beach community.

At its center is an emotionally charged relationship shaped by identity, loyalty, and the risks of first love.

Recommended for readers drawn to romantic YA fiction with strong atmosphere.'
	WHEN '10000000-0000-4000-8000-000000000094' THEN 'Eckhart Tolle invites readers to attend to the present moment and loosen identification with anxious thought and ego.

The approach is reflective and practical, offering concepts readers can test against their own routines and priorities.

Recommended for readers looking for a thoughtful starting point for personal change.'
	WHEN '10000000-0000-4000-8000-000000000095' THEN 'Through letters to an unseen friend, Charlie records friendship, grief, mental health, and the fragile discovery of belonging.

The story treats adolescence as a time of intense feeling, change, and connection, with themes that may be difficult for some readers.

Recommended for teen and adult readers who value emotionally candid coming-of-age fiction.'
	WHEN '10000000-0000-4000-8000-000000000096' THEN 'Madeline Miller reimagines the life of Patroclus and his bond with Achilles against the violence and glory of the Trojan War.

Period detail and social expectations frame the characters’ choices, adding context to the central emotional story.

Recommended for readers who enjoy historical settings, relationships, and moral dilemmas.'
	WHEN '10000000-0000-4000-8000-000000000097' THEN 'Neil Campbell’s textbook introduces the major systems, processes, and questions of modern biology through clear foundational concepts.

Definitions, examples, and structured explanations make it useful for study, review, and classroom discussion.

Recommended for students and general readers building a foundation in the subject.'
	WHEN '10000000-0000-4000-8000-000000000098' THEN 'Harry is unexpectedly entered in a dangerous international tournament as clues point toward a threat far beyond the competition.

Adventure, humor, and mythic stakes keep the story approachable while friendship and courage remain at its center.

Recommended for middle-grade readers and families who enjoy fast, imaginative series fiction.'
	WHEN '10000000-0000-4000-8000-000000000099' THEN 'A young girl’s apparent possession draws two priests into a crisis of faith, medicine, and terrifying supernatural possibility.

Tension builds through atmosphere, vulnerability, and the fear that familiar places or people may no longer be safe.

Recommended for readers who enjoy unsettling suspense and psychological horror.'
	WHEN '10000000-0000-4000-8000-000000000100' THEN 'A principled young woman and a reserved duke find their assumptions challenged in this Regency romance of family and social expectation.

Period detail and social expectations frame the characters’ choices, adding context to the central emotional story.

Recommended for readers who enjoy historical settings, relationships, and moral dilemmas.'
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
ELSE description
END;
