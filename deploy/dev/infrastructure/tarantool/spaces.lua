local s = box.schema.space.create('test_space')

s:format(
        {
            {name = 'id', type = 'integer'},
            {name = 'type_id', type = 'integer'},
            {name = 'status', type = 'integer'},
            {name = 'comment', type = 'string'},
        }
)

s:create_index('primary', {type = 'tree', parts = {'id', 'type_id'}})
s:create_index('id', {unique = false, type = 'tree', parts = {'id'}})

print("project-layout lua init test_space")