alter table configurations
    add column if not exists source_path text not null default '';

alter table configuration_revisions
    add column if not exists source_commit text not null default '',
    add column if not exists source_digest text not null default '';
