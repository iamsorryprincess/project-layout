#!/usr/bin/env tarantool

print('begin project-layout init.lua')

box.cfg {
    listen = '0.0.0.0:3301',
    pid_file = nil,
    background = false,
    log_level = 5,
    readahead = 32640,
    net_msg_max = 1536,
    memtx_memory = 1024 * 2 ^ 20,
    memtx_max_tuple_size = 8 * 1024 * 1024,
    slab_alloc_factor = 2.0
}

local function init()
    print('init from')
    box.schema.user.disable('guest')
    box.schema.user.create('test', {password = 'test', if_not_exists = true})
    box.schema.user.grant('test', 'read,write,execute', 'universe', nil, {if_not_exists = true})

    dofile('/usr/local/share/tarantool/spaces.lua')
    print('init spaces.lua')
end

box.once('init', init)