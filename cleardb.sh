mongodb='mongo 127.0.0.1:27017'
$mongodb <<EOF
use seele
db.dropDatabase()

exit;
EOF
