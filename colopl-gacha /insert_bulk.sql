USE gacha_system;

-- すでに同じマクロがある場合は一度削除する
DROP PROCEDURE IF EXISTS InsertBulkGachaLogs;

-- 10万件を高速自動生成するマクロ(ストアドプロシージャ)の定義
DELIMITER $$
CREATE PROCEDURE InsertBulkGachaLogs()
BEGIN
    DECLARE i INT DEFAULT 1;

    -- 【最適化】一括処理のためのトランザクション開始
    START TRANSACTION;

    -- 10万回ループを回す
    WHILE i <= 100000 DO
        INSERT INTO gacha_logs (user_id, item_id, rarity)
        VALUES (
            FLOOR(100 + (RAND() * 900)),                   -- 100~999のランダムなユーザーID（高カーディナリティの再現）
            FLOOR(1 + (RAND() * 1000)),                  -- 1~1000のランダムなアイテムID
            ELT(FLOOR(1 + (RAND() * 3)), 'SSR', 'SR', 'R') -- ランダムにレア度を割り振り
        );
        SET i = i + 1;
    END WHILE;

    -- 【最適化】一括コミットしてディスクへの書き込みを一回に抑える
    COMMIT;

END$$
DELIMITER ;

-- マクロを実際に呼び出して10万件のデータを挿入する
CALL InsertBulkGachaLogs();
