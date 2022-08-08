use sa_accounts;

# Simple Accounts

CREATE TABLE `sa_ac_type`
(
    `type`  varchar(10) NOT NULL COMMENT 'External value of account type',
    `value` smallint(6) NOT NULL COMMENT 'Internal value of account type',
    PRIMARY KEY (`type`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT ='Account type enumeration';
# 00000001011
INSERT INTO sa_ac_type (type, value)
VALUES ('ASSET', 11);
# 00000011011
INSERT INTO sa_ac_type (type, value)
VALUES ('BANK', 27);
# 00000000101
INSERT INTO sa_ac_type (type, value)
VALUES ('CR', 5);
# 00000101100
INSERT INTO sa_ac_type (type, value)
VALUES ('CUSTOMER', 44);
# 00000000011
INSERT INTO sa_ac_type (type, value)
VALUES ('DR', 3);
# 00000000000
INSERT INTO sa_ac_type (type, value)
VALUES ('DUMMY', 0);
# 01010000101
INSERT INTO sa_ac_type (type, value)
VALUES ('EQUITY', 645);
# 00001001101
INSERT INTO sa_ac_type (type, value)
VALUES ('EXPENSE', 77);
# 00110000101
INSERT INTO sa_ac_type (type, value)
VALUES ('INCOME', 389);
# 00010000101
INSERT INTO sa_ac_type (type, value)
VALUES ('LIABILITY', 133);
# 00000000001
INSERT INTO sa_ac_type (type, value)
VALUES ('REAL', 1);
# 10010000101
INSERT INTO sa_ac_type (type, value)
VALUES ('SUPPLIER', 1157);

CREATE TABLE `sa_coa`
(
    `id`   int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id of chart',
    `name` varchar(20)      NOT NULL COMMENT 'name of chart',
    PRIMARY KEY (`id`),
    UNIQUE KEY `sa_coa_name_uindex` (`name`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 7
  DEFAULT CHARSET = utf8 COMMENT ='A Chart of Account';

CREATE TABLE `sa_coa_ledger`
(
    `id`      int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'internal ledger id',
    `prntId`  int(10) unsigned NOT NULL DEFAULT 0 COMMENT 'parent node internal id',
    `lft`     int(10) unsigned NOT NULL DEFAULT 0 COMMENT 'left node internal id',
    `rgt`     int(10) unsigned NOT NULL DEFAULT 0 COMMENT 'node node internal id',
    `chartId` int(10) unsigned          DEFAULT NULL COMMENT 'id of chart that this account belongs to',
    `nominal` char(10)         NOT NULL COMMENT 'nominal id for this account',
    `type`    varchar(10)               DEFAULT NULL COMMENT 'type of account',
    `name`    varchar(30)      NOT NULL COMMENT 'name of account',
    `acDr`    bigint(20)       NOT NULL DEFAULT '0' COMMENT 'debit amount',
    `acCr`    bigint(20)       NOT NULL DEFAULT '0' COMMENT 'credit amount',
    PRIMARY KEY (`id`),
    UNIQUE KEY `sa_coa_ledger_chartId_nominal_index` (`chartId`, `nominal`),
    KEY `sa_coa_ledger_sa_ac_type_type_fk` (`type`),
    KEY `sa_coa_ledger_sa_coa_fk` (`chartId`),
    INDEX `sa_coa_ledger_lft_idx` (`lft`),
    INDEX `sa_coa_ledger_rgt_idx` (`rgt`),
    CONSTRAINT `sa_coa_ledger_sa_ac_type_type_fk` FOREIGN KEY (`type`) REFERENCES `sa_ac_type` (`type`) ON DELETE CASCADE,
    CONSTRAINT `sa_coa_ledger_sa_coa_fk` FOREIGN KEY (`chartId`) REFERENCES `sa_coa` (`id`) ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT ='Chart of Account structure';

CREATE TABLE `sa_journal`
(
    `id`      int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'internal id of the journal',
    `chartId` int(10) unsigned NOT NULL COMMENT 'the chart to which this journal belongs',
    `note`    text COMMENT 'a note for the journal entry',
    `date`    datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'timestamp for this journal',
    `src`     VARCHAR(6) COMMENT 'user defined source of journal',
    `ref`     INT(10) UNSIGNED COMMENT 'user defined reference to this journal',
    PRIMARY KEY (`id`),
    KEY `sa_journal_sa_coa_id_fk` (`chartId`),
    KEY `sa_journal_external_reference` (`src`, `ref`),
    CONSTRAINT `sa_journal_sa_coa_id_fk` FOREIGN KEY (`chartId`) REFERENCES `sa_coa` (`id`) ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT ='Txn Journal Header';

CREATE TABLE `sa_journal_entry`
(
    `id`      int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'internal id for entry',
    `jrnId`   int(10) unsigned DEFAULT NULL COMMENT 'id if journal that this entry belongs to',
    `nominal` varchar(10)      NOT NULL COMMENT 'nominal code for entry',
    `acDr`    bigint(20)       DEFAULT '0' COMMENT 'debit amount for entry',
    `acCr`    bigint(20)       DEFAULT '0' COMMENT 'credit amount for entry',
    PRIMARY KEY (`id`),
    KEY `sa_journal_entry_sa_org_id_fk` (`jrnId`),
    CONSTRAINT `sa_journal_entry_sa_jrn_id_fk` FOREIGN KEY (`jrnId`) REFERENCES `sa_journal` (`id`) ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT ='Txn Journal Entry';

CREATE
    DEFINER = CURRENT_USER FUNCTION
    sa_fu_add_chart(
    name VARCHAR(20)
)
    RETURNS INT(10) UNSIGNED
    MODIFIES SQL DATA DETERMINISTIC
BEGIN
    INSERT INTO sa_coa (`name`) VALUES (name);
    RETURN last_insert_id();
END;

CREATE
    DEFINER = CURRENT_USER PROCEDURE
    sa_sp_add_ledger(
    chartInternalId INT(10) UNSIGNED,
    nominal VARCHAR(10),
    type VARCHAR(10),
    name VARCHAR(30),
    prntNominal VARCHAR(10)
)
    MODIFIES SQL DATA DETERMINISTIC
BEGIN
    DECLARE vPrntId INT(10) UNSIGNED;
    DECLARE cntPrnts INT;
    DECLARE rightChildId INT;
    DECLARE myLeft INT;
    DECLARE myRight INT;

    # check to see if we already have a root account
    IF (prntNominal = '')
    THEN
        SELECT count(id)
        FROM sa_coa_ledger l
        WHERE l.prntId = 0
          AND l.chartId = chartInternalId
        INTO cntPrnts;

        IF (cntPrnts > 0)
        THEN
            SIGNAL SQLSTATE '45000'
                SET MYSQL_ERRNO = 1859, MESSAGE_TEXT = _utf8'Chart already has root account';
        END IF;
    END IF;

    SET vPrntId := 0;
    # Find the parent ledger id if the nominal id is not empty
    # as id cannot be zero, return zero if not found
    IF (prntNominal != '')
    THEN
        SELECT IFNULL((SELECT id
                       from sa_coa_ledger l
                       WHERE l.nominal = prntNominal
                         AND l.chartId = chartInternalId), 0)
        INTO vPrntId;

        IF (vPrntId = 0)
        THEN
            SIGNAL SQLSTATE '45000'
                SET MYSQL_ERRNO = 1107, MESSAGE_TEXT = _utf8'Invalid parent account nominal';
        END IF;
    END IF;

    IF (vPrntId = 0)
    THEN
        # We are inserting the root node - easy case
        INSERT INTO sa_coa_ledger (`prntId`, `lft`, `rgt`, `chartId`, `nominal`, `type`, `name`)
        VALUES (0, 1, 2, chartInternalId, nominal, type, name);
    ELSE
        # Does the parent have any children?
        SELECT IFNULL((SELECT max(id)
                       FROM sa_coa_ledger
                       WHERE prntId = vPrntId), 0)
        INTO rightChildId;

        IF (rightChildId = 0)
        THEN
            # no children
            SELECT lft
            FROM sa_coa_ledger
            WHERE id = vPrntId
            INTO myLeft;

            UPDATE sa_coa_ledger
            SET rgt = rgt + 2
            WHERE rgt > myLeft;

            UPDATE sa_coa_ledger
            SET lft = lft + 2
            WHERE lft > myLeft;

            INSERT INTO sa_coa_ledger (`prntId`, `lft`, `rgt`, `chartId`, `nominal`, `type`, `name`)
            VALUES (vPrntId, myLeft + 1, myLeft + 2, chartInternalId, nominal, type, name);
        ELSE
            # has children, add to right of last child
            SELECT rgt
            FROM sa_coa_ledger
            WHERE id = rightChildId
            INTO myRight;

            UPDATE sa_coa_ledger
            SET rgt = rgt + 2
            WHERE rgt > myRight;

            UPDATE sa_coa_ledger
            SET lft = lft + 2
            WHERE lft > myRight;

            INSERT INTO sa_coa_ledger (`prntId`, `lft`, `rgt`, `chartId`, `nominal`, `type`, `name`)
            VALUES (vPrntId, myRight + 1, myRight + 2, chartInternalId, nominal, type, name);
        END IF;
    END IF;
END;

CREATE
    DEFINER = CURRENT_USER PROCEDURE
    sa_sp_del_ledger(
    chartId INT(10) UNSIGNED,
    nominal VARCHAR(10)
)
    MODIFIES SQL DATA DETERMINISTIC
BEGIN
    DECLARE accId INT(10) UNSIGNED;
    DECLARE accDr INT(10) UNSIGNED;
    DECLARE accCr INT(10) UNSIGNED;
    SELECT id,
           acDr,
           acCr
    FROM sa_coa_ledger l
    WHERE l.nominal = nominal
      AND l.chartId = chartId
    INTO accId, accDr, accCr;

    IF (accDr > 0 OR accCr > 0)
    THEN
        SIGNAL SQLSTATE '45000'
            SET MYSQL_ERRNO = 2000, MESSAGE_TEXT = _utf8'Account balance is non zero';
    END IF;

    DELETE
    FROM sa_coa_ledger
    WHERE prntId = accId;

    DELETE
    FROM sa_coa_ledger
    WHERE id = accId;
END;

CREATE
    DEFINER = CURRENT_USER FUNCTION
    sa_fu_add_txn(
    chartId INT(10) UNSIGNED,
    note TEXT,
    date DATETIME,
    src VARCHAR(6),
    ref INT(10) UNSIGNED,
    arNominals TEXT,
    arAmounts TEXT,
    arTxnType TEXT
)
    RETURNS INT(10) UNSIGNED
    MODIFIES SQL DATA DETERMINISTIC
BEGIN
    DECLARE jrnId INT(10) UNSIGNED;
    DECLARE numInArray INT;

    SET date = IFNULL(date, CURRENT_TIMESTAMP);

    INSERT INTO sa_journal (`chartId`, `note`, `date`, `src`, `ref`)
    VALUES (chartId, note, date, src, ref);

    SELECT last_insert_id()
    INTO jrnId;

    SET numInArray =
                    char_length(arNominals) - char_length(replace(arNominals, ',', '')) + 1;

    SET @x = numInArray;
    REPEAT
        SET @txnType = substring_index(substring_index(arTxnType, ',', @x), ',', -1);
        SET @nominal = substring_index(substring_index(arNominals, ',', @x), ',', -1);
        SET @drAmount = 0;
        SET @crAmount = 0;
        IF @txnType = 'dr'
        THEN
            SET @drAmount = substring_index(substring_index(arAmounts, ',', @x), ',',
                                            -1);
        ELSE
            SET @crAmount = substring_index(substring_index(arAmounts, ',', @x), ',',
                                            -1);
        END IF;

        INSERT INTO sa_journal_entry (`jrnId`, `nominal`, `acDr`, `acCr`)
            VALUE (jrnId, @nominal, @drAmount, @crAmount);
        SET @x = @x - 1;
    UNTIL @x = 0 END REPEAT;

    RETURN jrnId;
END;

CREATE
    DEFINER = CURRENT_USER PROCEDURE
    sa_sp_get_tree(
    chartId INT(10) UNSIGNED
)
    READS SQL DATA
BEGIN
    SELECT prntId as origid,
           id     as destid,
           nominal,
           name,
           type,
           acDr,
           acCr
    FROM sa_coa_ledger
    WHERE `chartId` = chartId
    ORDER BY origid, destid;
END;

CREATE DEFINER = CURRENT_USER TRIGGER sp_tr_jrn_entry_updt
    AFTER INSERT
    ON sa_journal_entry
    FOR EACH ROW
BEGIN

    # get the internal ledger id
    SELECT l.id
    FROM sa_coa_ledger l
             LEFT JOIN sa_journal j ON j.chartId = l.chartId
    WHERE l.nominal = NEW.nominal
      AND j.id = NEW.jrnId
    INTO @acId;

    # create a concatenated string of parent ids as
    # creating temporary tables to hold parent ledger ids
    # borks, and you can't select from a table whilst updating

    SELECT GROUP_CONCAT(DISTINCT parent.id SEPARATOR ',')
    FROM sa_coa_ledger AS node,
         sa_coa_ledger AS parent
    WHERE node.lft BETWEEN parent.lft AND parent.rgt
      AND node.id = @acId
    GROUP BY node.id
    INTO @parents;

    SET @numInArray =
                    char_length(@parents) - char_length(replace(@parents, ',', '')) + 1;

    # update the parent ledger accounts
    WHILE (@numInArray > 0)
        DO
            SET @prntId = substring_index(substring_index(@parents, ',', @numInArray), ',', -1);
            UPDATE sa_coa_ledger l
            SET l.acDr = l.acDr + NEW.acDr,
                l.acCr = l.acCr + NEW.acCr
            WHERE id = @prntId;
            SET @numInArray = @numInArray - 1;
        END WHILE;
END;

# Application specific
CREATE TABLE `org_coa`
(
    `orgId` int(10) unsigned not null comment 'organisation id',
    `chartId` int(10) unsigned not null comment 'chart id',
    `chartUse` enum('CivilRecovery') not null default 'CivilRecovery' comment 'what the chart is used for',
    `crcy` enum('GBP','USD','EUR') not null default 'GBP' comment 'currency for chart',
    UNIQUE (`orgId`, `chartId`, `chartUse`),
    KEY `org_coa_chart_id_sa_coa` (`chartId`),
    CONSTRAINT `org_coa_chart_id_sa_coa_fk` FOREIGN KEY (`chartId`) REFERENCES `sa_coa` (`id`) ON DELETE RESTRICT
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT ='Organisation to COA links';

CREATE TABLE `org_cntrl`
(
    `chartId` int(10) unsigned not null comment 'chart id',
    `mnemonic` varchar(6) not null comment 'mnemonic to use for retrieval',
    `reason` varchar(30) comment 'reason for having the control account',
    `nominal` varchar(10) not null comment 'the nominal account id',
    UNIQUE (`chartId`, `mnemonic`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT ='Control account links to COA';



