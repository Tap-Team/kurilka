BEGIN;


create table if not exists subscription_types ( 
    id smallserial primary key,
    type varchar(16) not null,

    constraint subscription_types_unique unique (type)
);

create table if not exists motivations (
    id smallserial primary key,
    motivation text not null,

    constraint motivations_unique unique (motivation)
);

create table if not exists welcome_motivations (
    id smallserial primary key,
    motivation text not null,

    constraint welcome_motivations_unique unique (motivation)
);

CREATE OR REPLACE FUNCTION min_motivation_id() RETURNS smallint LANGUAGE SQL AS $$ 
    SELECT min(id) FROM motivations;
$$;

CREATE OR REPLACE FUNCTION min_welcome_motivation_id() RETURNS smallint LANGUAGE SQL AS $$ 
    SELECT min(id) FROM welcome_motivations;
$$;

create table if not exists users (
    id bigint not null,
    name varchar(15) not null,

    cigarette_day_amount int not null,
    cigarette_pack_amount int not null,
    pack_price real not null,
    abstinence_time timestamp(0) not null default now(),
    deleted boolean not null default FALSE,
    motivation_id smallint not null default min_motivation_id(),
    welcome_motivation_id smallint not null default min_welcome_motivation_id(),

    constraint users__motivations foreign key (motivation_id) references motivations(id),

    constraint users_welcome_motivations foreign key (welcome_motivation_id) references welcome_motivations(id),

    constraint users_key primary key (id)
);

create table if not exists user_subscriptions (
    user_id bigint not null,
    type_id bigint not null,

    expired timestamp(0),

    constraint user_subscriptions__users foreign key (user_id) references users(id) on delete cascade,
    constraint user_subscriptions__subscription_types foreign key (type_id) references subscription_types(id),

    constraint user_subscriptions_key primary key (user_id)
);


create table if not exists achievements_type (
    id smallserial primary key,
    type varchar(30) not null,

    constraint achievements_type_unique unique (type)
);

create table if not exists achievements (
    id smallserial primary key,
    level smallint not null,
    description text not null,
    motivation text not null,
    type_id smallint not null,
    exp integer not null,

    constraint achievements__achievements_type foreign key (type_id) references achievements_type(id),
    constraint achievements_description_unique unique (description),
    constraint achievements_unique unique (level, type_id)
);

create table if not exists user_achievements (
    achievement_id smallint not null,
    user_id bigint not null, 
    open_date timestamp(0),
    reach_date timestamp(0) default now(),
    shown boolean not null default FALSE,

    constraint user_achievements_key primary key (achievement_id, user_id),

    constraint user_achievements__users foreign key (user_id) references users(id) on delete cascade,
    constraint user_achievements__achievements foreign key (achievement_id) references achievements (id) on delete cascade
);

create table if not exists levels (
    level smallint not null,
    rank varchar(30) not null,
    min_exp integer not null,
    max_exp integer not null,

    constraint level_key primary key (level)
);

create table if not exists privacy_settings (
    id smallserial primary key,
    type varchar(30) not null,

    constraint privacy_settings_unique unique (type)
);

create table if not exists user_privacy_settings (
    user_id bigint not null,
    setting_id smallint not null,

    constraint user_privacy_settings_unique primary key (user_id, setting_id),

    constraint user_privacy_settings__users foreign key (user_id) references users(id) on delete cascade,

    constraint user_privacy_settings__privacy_settings foreign key (setting_id) references privacy_settings(id)
);


create table if not exists triggers (
    id smallserial primary key,
    name varchar(50) not null,

    constraint triggers_name_unique unique(name)
);

create table if not exists user_triggers (
    trigger_id smallint not null,
    user_id bigint not null,

    constraint user_triggers__users foreign key (user_id) references users(id) on delete cascade,
    constraint user_triggers__triggers foreign key (trigger_id) references triggers(id),
    constraint user_triggers_primary primary key (trigger_id, user_id)
);


COMMIT;


BEGIN;

INSERT INTO privacy_settings (type) VALUES 
('STATISTICS_MONEY'),
('STATISTICS_CIGARETTE'),
('STATISTICS_LIFE'),
('STATISTICS_TIME'),
('ACHIEVEMENTS_DURATION'),
('ACHIEVEMENTS_HEALTH'),
('ACHIEVEMENTS_WELL_BEING'),
('ACHIEVEMENTS_SAVING'),
('ACHIEVEMENTS_CIGARETTE');

INSERT INTO triggers (name) VALUES 
('THANK_YOU'),
('SUPPORT_CIGGARETTE'),
('SUPPORT_HEALTH'),
('SUPPORT_TRIAL');


INSERT INTO subscription_types (type) VALUES ('NONE'), ('TRIAL'), ('BASIC');

INSERT INTO levels (level, rank, min_exp, max_exp) VALUES 
(1, 'Новичок', 0, 99),
(2, 'Новичок', 100, 199),
(3, 'Опытный', 200, 299),
(4, 'Опытный', 300, 399),
(5, 'Уверенный', 400, 499),
(6, 'Уверенный', 500, 599),
(7, 'Бывалый', 600, 699),
(8, 'Бывылый', 700, 799),
(9, 'Профессионал', 800, 899),
(10, 'Мастер', 900, 1000);


COMMIT;




BEGIN;


CREATE OR REPLACE FUNCTION insert_motivations(motivations_array text[])
RETURNS void AS $$
DECLARE
    motivation_text text;
BEGIN
    FOREACH motivation_text IN ARRAY motivations_array
    LOOP
        INSERT INTO motivations (motivation) VALUES (motivation_text);
    END LOOP;
END;
$$ LANGUAGE plpgsql;


SELECT insert_motivations(
    ARRAY[
        'Первые 3 дня самые сложные',
        'Сигареты ничего не дают. Это трата времени, денег и здоровья.',
        'Думайте о курении как о болезни. Вы идёте на поправку и совсем скоро симптомы, которые призывают вас покурить - пройдут.',
        'Дайте себе обещание в этот день не курить ни при каких обстоятельствах. Завтра дайте его ещё раз.',
        'Новая жизнь без сигарет - лучше старой! Не возвращайтесь к ним никогда!',
        'Чтобы бросить курить - нужно просто не тратить время на сигареты. Лень - наш союзник!',
        'Избавьтесь от всех вещей и атрибутов, которые напоминают Вам о курении.',
        'Нельзя курить даже одну сигарету! Это ловушка, чтобы снова погрузиться в мир курения и несчастий.',
        'Подумайте о плюсах, которыми вы обладаете благодаря избавлению от пагубной привычки!',
        'Сэкономленные деньги от сигарет можно потратить на полезные вещи или развлечения.',
        'Больше не будет неприятного запаха изо рта, от рук и от одежды, а также в квартире!',
        'Курение может стать причиной многих болезней.',
        'Не сдавайся! Ты можешь победить зависимость!', 
        'Бросить курить - это выигрыш, а не потеря.',
        'Самое время жить полноценную жизнь без курения.',
        'Курение лишь крадёт у тебя время и здоровье.',
        'Помни о том, что курение приносит только негативные последствия.',
        'Помни, что каждый день без сигарет - это шаг к победе над зависимостью.',
        'Бывшие курильщики живут дольше чем те, кто всё ещё продолжает курить.',
        'Прекращение курение снижает риск развития рака лёгких, болезней сердца, инсультов, хронических заболеваний лёгких.',
        'Очень скоро твоё настроение перестанет зависеть от того, покурил ты или нет.',
        'Избавьтесь от всех вещей, которые напоминали вам о курении.',
        'В первые дни отказа от курения помогут справиться леденцы.' ,
        'Средний срок жизни курильщика меньше на 10-15 лет.',
        'Никотин разрушает клетки мозга, что снижает умственные способности и ухудшает память.',
        'Уже в 40-45 лет курильщик чувствует себя так, как некурящий человек в 55-60 лет.', 
        'Попроси семью и друзей поддержать тебя в твоём решении бросить курить.', 
        'Старайся не пить крепкий чай или кофе – это обостряет тягу к сигаретам.',
        'Физкультура не только отвлекает от курения, но и очищает дыхание.',
        'Поставьте перед собой цель. Имейте твёрдое желания вылечиться от пагубной привычки.'
    ]
);


COMMIT;


BEGIN;


CREATE OR REPLACE FUNCTION insert_welcome_motivations(motivations_array text[])
RETURNS void AS $$
DECLARE
    motivation_text text;
BEGIN
    FOREACH motivation_text IN ARRAY motivations_array
    LOOP
        INSERT INTO welcome_motivations (motivation) VALUES (motivation_text);
    END LOOP;
END;
$$ LANGUAGE plpgsql;

SELECT insert_welcome_motivations(
    ARRAY[
        'Ты сильнее, чем сигарета!',
        'Не сдавайся!',
        'Будь здоровым!',
        'Побеждай себя и дыши полной грудью!',
        'Продолжай заботиться о себе!',
        'Жизнь без дыма и никотина - класс!',
        'Терпение приведет к успеху!',
        'Мы в тебя верим!',
        'Жизнь без сигарет лучше прежний!',
        'Стань ещё лучше сегодня!',
        'Здоровье с каждым днем лучше!',
        'Наслаждайся жизнью!',
        'Начни инвестировать в здоровье!',
        'Нет курения - больше энергии!',
        'Меньше дыма - больше жизни!',
        'Достигай успеха без курения!',
        'Борьба с курением - путь к здоровью!',
        'Ты можешь преодолеть это!',
        'Брось курить, чтобы жить свободно!',
        'Заботься о своем теле!'
    ]
);

COMMIT;


BEGIN;

INSERT INTO achievements_type (type) VALUES 
('Длительность'),
('Сигареты'),
('Здоровье'),
('Самочувствие'),
('Экономия');


CREATE OR REPLACE FUNCTION insert_achievement(ach_type varchar(30), ach_level int, ach_exp int, ach_description text, ach_motivation text)
RETURNS void AS $$
BEGIN
    INSERT INTO achievements (type_id, level, exp, description,motivation) VALUES (
        (SELECT id FROM achievements_type WHERE type = ach_type),
        ach_level,
        ach_exp,
        ach_description,
        ach_motivation
    );
END;
$$ LANGUAGE plpgsql;


SELECT insert_achievement('Длительность',1,20,'Вы не курили 1 день - 20 xp','Твоя жизнь станет ярче и здоровее без курения.');
SELECT insert_achievement('Длительность',2,20,'Вы не курили 3 дня - 20 xp','С каждым днем без сигарет ты становишься сильнее.');
SELECT insert_achievement('Длительность',3,20,'Вы не курили 1 нед - 20 xp','Твоя решимость и терпение наградят тебя здоровьем и счастьем.');
SELECT insert_achievement('Длительность',4,20,'Вы не курили 1 мес - 20 xp','Продолжай бить свои собственные рекорды дней без сигарет.');
SELECT insert_achievement('Длительность',5,20,'Вы не курили 2 мес - 20 xp','Новый здоровый образ жизни уже ждет тебя с открытыми объятиями.');
SELECT insert_achievement('Длительность',6,20,'Вы не курили 3 мес - 20 xp','Будь сильным и настойчивым, ты отлично справляешься.');
SELECT insert_achievement('Длительность',7,20,'Вы не курили 6 мес - 20 xp','Ты уже показал, что можешь контролировать свою зависимость, продолжай этот путь.');
SELECT insert_achievement('Длительность',8,20,'Вы не курили 9 мес - 20 xp','Никакая сигарета не стоит того, чтобы пожертвовать своим здоровьем.');
SELECT insert_achievement('Длительность',9,20,'Вы не курили 1 год - 20 xp','Жизнь после курения ярче и полна новых впечатлений.');
SELECT insert_achievement('Длительность',10,20,'Вы не курили 1.5 года - 20 xp','Не позволяй никакому желанию закурить загубить твой успех.');


SELECT insert_achievement('Сигареты',1,20,'Вы не выкурили 20 сиг - 20 xp','Ты на верном пути к здоровой и счастливой жизни.');
SELECT insert_achievement('Сигареты',2,20,'Вы не выкурили 50 сиг - 20 xp','Будь настойчивым и терпеливым - твои результаты придут со временем.');
SELECT insert_achievement('Сигареты',3,20,'Вы не выкурили 100 сиг - 20 xp','Твоя решимость и сила воли помогут тебе преодолеть этот вызов.');
SELECT insert_achievement('Сигареты',4,20,'Вы не выкурили 250 сиг - 20 xp','Награда за твои усилия уже близко.');
SELECT insert_achievement('Сигареты',5,20,'Вы не выкурили 500 сиг - 20 xp','Помни, что твое здоровье - твой самый ценный актив.');
SELECT insert_achievement('Сигареты',6,20,'Вы не выкурили 750 сиг - 20 xp','Никакое желание закурить не может остановить твою решимость.');
SELECT insert_achievement('Сигареты',7,20,'Вы не выкурили 1000 сиг - 20 xp','Каждый день без сигарет - это новый шанс стать лучше.');
SELECT insert_achievement('Сигареты',8,20,'Вы не выкурили 1500 сиг - 20 xp','Продолжай двигаться вперед без сигарет.');
SELECT insert_achievement('Сигареты',9,20,'Вы не выкурили 2000 сиг - 20 xp','Наслаждайся свободой от курения.');
SELECT insert_achievement('Сигареты',10,20,'Вы не выкурили 3000 сиг - 20 xp','Твои близкие и друзья гордятся тобой за твою силу воли.');


SELECT insert_achievement('Здоровье',1,20,'Ваш пульс снова в норме (20 мин) - 20 xp','Будь героем своей жизни и брось курить.');
SELECT insert_achievement('Здоровье',2,20,'Риск сердечного приступа начинает снижаться (8 ч) - 20 xp','Продолжай двигаться вперед без сигарет.');
SELECT insert_achievement('Здоровье',3,20,'Обогащение крови кислородом снова в норме (9 ч) - 20 xp','Наслаждайся своим новым здоровым образом жизни.');
SELECT insert_achievement('Здоровье',4,20,'Угарный газ полностью выведен из организма (24 ч) - 20 xp','Каждый день, побеждая желание закурить, ты делаешь себя сильнее.');
SELECT insert_achievement('Здоровье',5,20,'Появился кашель, но это нормально. Ваш организм очищается (33 ч) - 20 xp','Не забывай, что каждый день без сигарет - это победа.');
SELECT insert_achievement('Здоровье',6,20,'В Вашей крови больше не осталось никотина (48 ч) - 20 xp','Твои легкие уже счастливы без сигарет.');
SELECT insert_achievement('Здоровье',7,20,'Кашель и нагрузка на бронхи  сходит на нет (3 нед.) - 20 xp','Не позволяй никакому желанию закурить загубить твой успех.');
SELECT insert_achievement('Здоровье',8,20,'Снижение риска респираторных заболеваний (1 мес) - 20 xp','Никакая сигарета не стоит того, чтобы пожертвовать своим здоровьем.');
SELECT insert_achievement('Здоровье',9,20,'Эффективность работы Ваших лёгких выросла примерно на 10% (3 мес) - 20 xp','Ты заслуживаешь здоровой, счастливой жизни без курения.');
SELECT insert_achievement('Здоровье',10,20,'Уровень Ваш иммунитет улучшился (6 мес) - 20 xp','Жизнь без курения полна новых возможностей и перспектив.');


SELECT insert_achievement('Экономия',1,20,'Сэкономлено 300 р - 20 xp','Ты уже сделал первый шаг к здоровой и счастливой жизни.');
SELECT insert_achievement('Экономия',2,20,'Сэкономлено 600 р - 20 xp','Награда за твои усилия будет значительнее, чем ты думаешь.');
SELECT insert_achievement('Экономия',3,20,'Сэкономлено 1 200 р - 20 xp','Твои близкие и друзья поддерживают тебя на этом пути.');
SELECT insert_achievement('Экономия',4,20,'Сэкономлено 2 000 р - 20 xp','Не дай желанию закурить затмить твои успехи.');
SELECT insert_achievement('Экономия',5,20,'Сэкономлено 3 000 р - 20 xp','Сэкономленные деньги можно потратить на новые цели.');
SELECT insert_achievement('Экономия',6,20,'Сэкономлено 4 000 р - 20 xp','Наслаждайся свободой от курения и новыми возможностями, которые тебе ждут.');
SELECT insert_achievement('Экономия',7,20,'Сэкономлено 6 000 р - 20 xp','Твоя решимость и настойчивость помогут тебе победить желание закурить.');
SELECT insert_achievement('Экономия',8,20,'Сэкономлено 10 000 р - 20 xp','Каждый день без сигарет - это новая победа над своей зависимостью.');
SELECT insert_achievement('Экономия',9,20,'Сэкономлено 15 000 р - 20 xp','Ты делаешь что-то действительно важное для своего здоровья и жизни.');
SELECT insert_achievement('Экономия',10,20,'Сэкономлено 25 000р - 20 xp','Ты уже доказал свою силу воли и решимость, продолжай двигаться вперед.');


SELECT insert_achievement('Самочувствие',1,20,'К вам вернулось чувство вкуса и пища стала казаться вкуснее (5 дней) - 20 xp','Твое решение бросить курить - это самый лучший подарок для твоего здоровья.');
SELECT insert_achievement('Самочувствие',2,20,'Ваше обоняние обострилось (7 дней) - 20 xp','Ты можешь преодолеть любые трудности, включая желание закурить.');
SELECT insert_achievement('Самочувствие',3,20,'Физические нагрузки и даже дыхание стали даваться легче (10 дней) - 20 xp','Будь настойчив и терпелив, результаты твоих усилий скоро появятся.');
SELECT insert_achievement('Самочувствие',4,20,'Ваша кожа приобрела более светлый оттенок (14 дней)- 20 xp','Ты на верном пути к своей лучшей жизни.');
SELECT insert_achievement('Самочувствие',5,20,'Ваша раздражительность снизилась, а качество сна значительно улучшилось (17 дней) - 20 xp','Продолжай идти вперед, став победителем в этой битве.');
SELECT insert_achievement('Самочувствие',6,20,'Ваш голос становится немного выше (25 дней) - 20 xp','Каждый день без сигарет - это новая жизнь и новый шанс стать лучшим.');
SELECT insert_achievement('Самочувствие',7,20,'Вы ощущаете свободу и большую уверенность в себе (31 дней) - 20 xp','Помни, что твое здоровье - твой самый ценный актив.');
SELECT insert_achievement('Самочувствие',8,20,'Ваше лицо избавилось от серости и угрюмости (35 дней) - 20 xp','Жизнь после курения ярче и полна новых впечатлений.');
SELECT insert_achievement('Самочувствие',9,20,'К Вам вернулось ваше либидо (40 дней) - 20 xp','Никакая сигарета не стоит того, чтобы подорвать твои усилия.');
SELECT insert_achievement('Самочувствие',10,20,'Вы стали меньше уставать (50 дней) - 20 xp','Наслаждайся своим новым, здоровым образом жизни.');

COMMIT;