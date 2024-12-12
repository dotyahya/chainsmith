import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import Home from './Home';
import BlockExplorer from './BlockExplorer';
import Transaction from './Transaction';
import Wallet from './Wallet';

function App() {
  return (
    <Router>
      <div className="min-h-screen bg-gray-100">
        {/* Navigation Bar */}
        <nav className="bg-blue-600 p-4">
          <ul className="flex space-x-4 text-white">
            <li>
              <Link to="/" className="hover:underline">Home</Link>
            </li>
            <li>
              <Link to="/block-explorer" className="hover:underline">Block Explorer</Link>
            </li>
            <li>
              <Link to="/transaction" className="hover:underline">Transactions</Link>
            </li>
            <li>
              <Link to="/wallet" className="hover:underline">Wallet</Link>
            </li>
          </ul>
        </nav>

        {/* Define routes for each screen */}
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/block-explorer" element={<BlockExplorer />} />
          <Route path="/transaction" element={<Transaction />} />
          <Route path="/wallet" element={<Wallet />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
