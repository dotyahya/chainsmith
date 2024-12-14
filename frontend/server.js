const express = require('express');
const cors = require('cors');
const bodyParser = require('body-parser');

const app = express();
app.use(cors());
app.use(bodyParser.json());

const blocks = [
  { id: 1, timestamp: '2024-12-14', hash: '0x123abc...' },
  { id: 2, timestamp: '2024-12-14', hash: '0x456def...' },
];

app.get('/api/blocks', (req, res) => {
  res.json(blocks);
});

const PORT = 3001;
app.listen(PORT, () => console.log(`Backend running on http://localhost:${PORT}`));
