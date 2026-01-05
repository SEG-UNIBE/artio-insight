SELECT pg_size_pretty( pg_total_relation_size('tablename') );

select pg_size_pretty(pg_total_relation_size(tablename)), * from pg_tables where schemaname = 'rdv'
show SEARCH_PATH
SET SEARCH_PATH = inb,rdv

select
    table_name,
    pg_size_pretty(pg_total_relation_size(quote_ident(table_name))),
    pg_total_relation_size(quote_ident(table_name))
from information_schema.tables
where table_schema = 'rdv'
order by 3 desc;

select
   pg_size_pretty(sum( pg_total_relation_size(quote_ident(table_name))))
from information_schema.tables
where table_schema = 'rdv'
order by 3 desc;



select count(*) from logs_artiostr_sat
union all
select count(*) from events_artiostr_sat_cur
union all
select count(*) from events_artiostr_sat
union all
select count(*) from logs_artioinf_sat_cur
union all
select count(*) from logs_artiostr_sat_cur
union all
select count(*) from logs_artioinf_sat
union all
select count(*) from dump_neo4j_sat
union all
select count(*) from dump_neo4j_sat_cur
union all
select count(*) from events_artioinf_sat_cur
union all
select count(*) from events_artioinf_sat
union all
select count(*) from nipkinds_nipfile_sat_cur
union all
select count(*) from nipkinds_nipfile_sat

select * from udm.neo4j_nodes where _labels = ':Relay'

select * from udm.software


SELECT *
FROM udm.relaysoftware AS rs
         INNER JOIN udm.software s
                   ON s.software_id = rs.software_id


SELECT s.software_id, s.software, count(r.relay_id)
FROM udm.relay AS r
         LEFT JOIN udm.relaysoftware AS rs ON rs.relay_id = r.relay_id
         LEFT JOIN udm.software AS s ON s.software_id = rs.software_id

group by s.software, s.software_id
order by 3 desc

SELECT *
FROM udm.relay AS r
         LEFT JOIN udm.relaysoftware AS rs ON rs.relay_id = r.relay_id
         LEFT JOIN udm.software AS s ON s.software_id = rs.software_id

where s.software = 'N/A' and relay_name ~ E'\n';