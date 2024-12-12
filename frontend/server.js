const express = require('express');
const app = express();
const port = 5000;

app.use(express.json());

// Dummy data
const dummyData = {
  status: { blockCount: 123, transactionCount: 456 },
  blocks: [{ id: 1, transactions: [{}, {}] }, { id: 2, transactions: [{}] }],
  transactions: [{ sender: 'Alice', receiver: 'Bob', amount: 10 }],
  wallet: { balance: 100, transactions: [{ amount: 10, date: '2024-12-12' }] },
};

app.get('/api/status', (req, res) => res.json(dummyData.status));
app.get('/api/blocks', (req, res) => res.json({ blocks: dummyData.blocks }));
app.get('/api/transactions', (req, res) => res.json({ transactions: dummyData.transactions }));
app.get('/api/wallet', (req, res) => res.json(dummyData.wallet));

app.post('/api/send-transaction', (req, res) => {
  console.log(req.body); // Transaction details
  res.json({ message: 'Transaction sent successfully!' });
});

app.listen(port, () => {
  console.log(`Server running at http://localhost:${port}`);
});
