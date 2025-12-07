import React, { useState, useEffect, useCallback } from 'react';
import { useAuth } from '../context/AuthContext';
import api from '../services/api';
import Navbar from '../components/Navbar';

const Mining = () => {
  const { token } = useAuth();
  const [miningStatus, setMiningStatus] = useState(null);
  const [blockchainStats, setBlockchainStats] = useState(null);
  const [myBlocks, setMyBlocks] = useState([]);
  const [recentBlocks, setRecentBlocks] = useState([]);
  const [isMining, setIsMining] = useState(false);
  const [miningResult, setMiningResult] = useState(null);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(true);
  const [showCoinbaseModal, setShowCoinbaseModal] = useState(false);
  const [coinbaseAmount, setCoinbaseAmount] = useState('');
  const [creatingCoinbase, setCreatingCoinbase] = useState(false);
  const [wallet, setWallet] = useState(null);

  const fetchData = useCallback(async () => {
    try {
      const [statusRes, statsRes, blocksRes, myBlocksRes, walletRes] = await Promise.all([
        api.blockchain.getMiningStatus(),
        api.blockchain.getStats(),
        api.blockchain.getBlocks(1, 5),
        api.blockchain.getMyBlocks(token),
        api.wallet.getMyWallet()
      ]);

      setMiningStatus(statusRes.data);
      setBlockchainStats(statsRes.data.stats);
      setRecentBlocks(blocksRes.data.blocks || []);
      setMyBlocks(myBlocksRes.data.blocks || []);
      setWallet(walletRes.data.wallet);
      setError('');
    } catch (err) {
      if (err.response?.status !== 404) {
        setError('Failed to load blockchain data');
      }
    } finally {
      setLoading(false);
    }
  }, [token]);

  useEffect(() => {
    if (token) {
      fetchData();
    }
  }, [token, fetchData]);

  const handleCreateCoinbase = async () => {
    if (!coinbaseAmount || isNaN(parseFloat(coinbaseAmount)) || parseFloat(coinbaseAmount) <= 0) {
      setError('Please enter a valid amount');
      return;
    }

    setCreatingCoinbase(true);
    setError('');

    try {
      await api.utxo.createCoinbase({
        walletId: wallet?.walletId,
        amount: parseFloat(coinbaseAmount),
        reason: 'initial_distribution'
      });
      setShowCoinbaseModal(false);
      setCoinbaseAmount('');
      fetchData();
      setMiningResult({
        type: 'coinbase',
        message: `Successfully created ${coinbaseAmount} coins!`
      });
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to create coinbase');
    } finally {
      setCreatingCoinbase(false);
    }
  };

  const handleCreateGenesis = async () => {
    try {
      setError('');
      const response = await api.blockchain.createGenesis(token);
      setMiningResult({
        type: 'genesis',
        message: response.data.message,
        block: response.data.block
      });
      fetchData();
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to create genesis block');
    }
  };

  const handleMine = async () => {
    setIsMining(true);
    setMiningResult(null);
    setError('');

    try {
      const response = await api.blockchain.mine(token);
      setMiningResult({
        type: 'mined',
        message: response.data.message,
        block: response.data.block,
        reward: response.data.miningReward,
        nonce: response.data.nonce
      });
      fetchData();
    } catch (err) {
      setError(err.response?.data?.error || 'Mining failed');
    } finally {
      setIsMining(false);
    }
  };

  const handleValidateChain = async () => {
    try {
      const response = await api.blockchain.validate();
      const validation = response.data.validation;
      if (validation.isValid) {
        alert(`✅ Blockchain is valid! ${validation.blocksChecked} blocks verified.`);
      } else {
        alert(`❌ Blockchain validation failed!\nErrors: ${validation.errors.join('\n')}`);
      }
    } catch (err) {
      setError('Failed to validate blockchain');
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen">
        <Navbar />
        <div className="flex items-center justify-center h-96">
          <div className="text-center">
            <div className="w-16 h-16 border-4 border-purple-500 border-t-transparent rounded-full animate-spin mx-auto"></div>
            <p className="text-gray-600 mt-4">Loading blockchain data...</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen">
      <Navbar />
      
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-800 mb-2">⛏️ Mining Center</h1>
          <p className="text-gray-600">Mine blocks and earn rewards with Proof-of-Work</p>
        </div>

        {error && (
          <div className="mb-6 bg-red-500/10 border border-red-500/50 text-red-400 px-4 py-3 rounded-xl flex items-center gap-2">
            <svg className="w-5 h-5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
            </svg>
            {error}
          </div>
        )}

        {miningResult && (
          <div className="mb-6 bg-green-500/10 border border-green-500/50 text-green-400 px-6 py-4 rounded-2xl">
            <div className="flex items-start gap-4">
              <div className="w-12 h-12 rounded-xl bg-green-500/20 flex items-center justify-center text-2xl flex-shrink-0">
                {miningResult.type === 'genesis' ? '🌟' : '🎉'}
              </div>
              <div>
                <h3 className="text-lg font-bold">
                  {miningResult.type === 'genesis' ? 'Genesis Block Created!' : 'Block Mined Successfully!'}
                </h3>
                <p className="text-green-300/80 mt-1">{miningResult.message}</p>
                {miningResult.reward && (
                  <p className="mt-2 text-lg">
                    Reward: <span className="font-bold text-green-400">{miningResult.reward} coins</span>
                  </p>
                )}
                {miningResult.nonce && (
                  <p className="text-sm text-green-300/60 mt-1">Nonce: {miningResult.nonce}</p>
                )}
              </div>
            </div>
          </div>
        )}

        {/* Stats Grid */}
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
          <div className="glass-card rounded-2xl p-6 card-hover">
            <div className="flex items-center justify-between mb-2">
              <span className="text-gray-600 text-sm">Total Blocks</span>
              <span className="text-2xl">📦</span>
            </div>
            <p className="text-3xl font-bold text-gray-800">{blockchainStats?.totalBlocks || 0}</p>
          </div>
          
          <div className="glass-card rounded-2xl p-6 card-hover">
            <div className="flex items-center justify-between mb-2">
              <span className="text-gray-600 text-sm">Pending Txs</span>
              <span className="text-2xl">⏳</span>
            </div>
            <p className="text-3xl font-bold text-yellow-600">{miningStatus?.pendingTransactions || 0}</p>
          </div>
          
          <div className="glass-card rounded-2xl p-6 card-hover">
            <div className="flex items-center justify-between mb-2">
              <span className="text-gray-600 text-sm">Difficulty</span>
              <span className="text-2xl">🎯</span>
            </div>
            <p className="text-3xl font-bold text-purple-600">{miningStatus?.difficulty || 4}</p>
            <p className="text-xs text-gray-500 mt-1">Target: {miningStatus?.targetPrefix || '0000'}</p>
          </div>
          
          <div className="glass-card rounded-2xl p-6 card-hover">
            <div className="flex items-center justify-between mb-2">
              <span className="text-gray-600 text-sm">Block Reward</span>
              <span className="text-2xl">💰</span>
            </div>
            <p className="text-3xl font-bold text-green-600">{miningStatus?.miningReward || 50}</p>
          </div>
        </div>

        {/* Mining Actions */}
        <div className="glass-card rounded-2xl p-6 mb-8">
          <h2 className="text-xl font-bold text-gray-800 mb-6">Mining Controls</h2>
          
          <div className="flex flex-wrap gap-4">
            {(!blockchainStats || blockchainStats.totalBlocks === 0) && (
              <button
                onClick={handleCreateGenesis}
                className="px-6 py-3 rounded-xl font-semibold text-black transition-all duration-300 transform hover:scale-105 flex items-center gap-2"
                style={{ background: 'linear-gradient(135deg, #f2c94c 0%, #f2994a 100%)' }}
              >
                <span className="text-xl">🌟</span>
                Create Genesis Block
              </button>
            )}
            
            <button
              onClick={handleMine}
              disabled={isMining || !blockchainStats || blockchainStats.totalBlocks === 0}
              className="btn-primary flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
            >
              {isMining ? (
                <>
                  <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                  </svg>
                  Mining...
                </>
              ) : (
                <>
                  <span className="text-xl">⛏️</span>
                  Mine Block
                </>
              )}
            </button>

            <button
              onClick={handleValidateChain}
              className="px-6 py-3 rounded-xl font-semibold text-white transition-all duration-300 transform hover:scale-105 flex items-center gap-2"
              style={{ background: 'linear-gradient(135deg, #11998e 0%, #38ef7d 100%)' }}
            >
              <span className="text-xl">✅</span>
              Validate Chain
            </button>

            {wallet && (
              <button
                onClick={() => setShowCoinbaseModal(true)}
                className="px-6 py-3 rounded-xl font-semibold text-white transition-all duration-300 transform hover:scale-105 flex items-center gap-2"
                style={{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }}
              >
                <span className="text-xl">💎</span>
                Create Coins
              </button>
            )}
          </div>

          {(miningStatus?.pendingTransactions || 0) === 0 && blockchainStats?.totalBlocks > 0 && (
            <div className="mt-4 p-4 bg-yellow-500/10 border border-yellow-500/30 rounded-xl flex items-center gap-3">
              <span className="text-2xl">💡</span>
              <p className="text-yellow-600">No pending transactions. You can still mine empty blocks for rewards!</p>
            </div>
          )}
        </div>

        {/* Coinbase Creation Modal */}
        {showCoinbaseModal && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
            <div className="bg-white rounded-2xl p-6 max-w-md w-full shadow-xl">
              <div className="flex items-center justify-between mb-6">
                <h3 className="text-xl font-bold text-gray-800">💎 Create Coins</h3>
                <button
                  onClick={() => setShowCoinbaseModal(false)}
                  className="text-gray-400 hover:text-gray-600"
                >
                  <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>

              <p className="text-gray-600 text-sm mb-4">
                Create new coins directly to your wallet. This is for testing purposes.
              </p>

              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Amount
                </label>
                <input
                  type="number"
                  value={coinbaseAmount}
                  onChange={(e) => setCoinbaseAmount(e.target.value)}
                  className="input-modern"
                  placeholder="Enter amount"
                  min="1"
                  step="0.01"
                />
              </div>

              <div className="bg-gray-50 rounded-xl p-4 mb-6">
                <p className="text-gray-600 text-sm">
                  <span className="font-medium">Recipient:</span> {wallet?.walletId?.substring(0, 20)}...
                </p>
              </div>

              <div className="flex gap-3">
                <button
                  onClick={() => setShowCoinbaseModal(false)}
                  className="flex-1 px-4 py-3 bg-gray-100 text-gray-700 rounded-xl font-semibold hover:bg-gray-200 transition-colors"
                >
                  Cancel
                </button>
                <button
                  onClick={handleCreateCoinbase}
                  disabled={creatingCoinbase || !coinbaseAmount}
                  className="flex-1 btn-primary disabled:opacity-50"
                >
                  {creatingCoinbase ? 'Creating...' : 'Create Coins'}
                </button>
              </div>
            </div>
          </div>
        )}

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* My Mining Stats */}
          <div className="glass-card rounded-2xl p-6">
            <h2 className="text-xl font-bold text-gray-800 mb-6">📊 My Mining Stats</h2>
            
            <div className="grid grid-cols-3 gap-4">
              <div className="bg-gray-50 rounded-xl p-4 text-center">
                <p className="text-gray-600 text-sm">Blocks Mined</p>
                <p className="text-2xl font-bold text-gray-800 mt-1">{myBlocks.length}</p>
              </div>
              <div className="bg-gray-50 rounded-xl p-4 text-center">
                <p className="text-gray-600 text-sm">Total Earned</p>
                <p className="text-2xl font-bold text-green-600 mt-1">
                  {myBlocks.reduce((sum, b) => sum + (b.miningReward || 0), 0).toFixed(0)}
                </p>
              </div>
              <div className="bg-gray-50 rounded-xl p-4 text-center">
                <p className="text-gray-600 text-sm">Last Block</p>
                <p className="text-2xl font-bold text-purple-600 mt-1">
                  {myBlocks.length > 0 ? `#${myBlocks[0]?.index}` : '-'}
                </p>
              </div>
            </div>
          </div>

          {/* Mining Info */}
          <div className="glass-card rounded-2xl p-6">
            <h2 className="text-xl font-bold text-gray-800 mb-6">ℹ️ How Mining Works</h2>
            
            <div className="space-y-4">
              <div className="flex gap-3">
                <div className="w-8 h-8 rounded-lg bg-purple-500/20 flex items-center justify-center text-purple-600 font-semibold flex-shrink-0">1</div>
                <p className="text-gray-700">Collect pending transactions from the mempool</p>
              </div>
              <div className="flex gap-3">
                <div className="w-8 h-8 rounded-lg bg-purple-500/20 flex items-center justify-center text-purple-600 font-semibold flex-shrink-0">2</div>
                <p className="text-gray-700">Find a nonce that produces a hash with {miningStatus?.difficulty || 4} leading zeros</p>
              </div>
              <div className="flex gap-3">
                <div className="w-8 h-8 rounded-lg bg-purple-500/20 flex items-center justify-center text-purple-600 font-semibold flex-shrink-0">3</div>
                <p className="text-gray-700">Submit the valid block and receive {miningStatus?.miningReward || 50} coins reward</p>
              </div>
            </div>
          </div>
        </div>

        {/* Recent Blocks */}
        <div className="mt-6 glass-card rounded-2xl p-6">
          <h2 className="text-xl font-bold text-gray-800 mb-6">🔗 Recent Blocks</h2>
          
          {recentBlocks.length === 0 ? (
            <div className="text-center py-12">
              <div className="w-16 h-16 bg-gray-100 rounded-2xl flex items-center justify-center mx-auto mb-4">
                <span className="text-3xl">📦</span>
              </div>
              <h3 className="text-gray-800 font-semibold mb-2">No Blocks Yet</h3>
              <p className="text-gray-600">Create the genesis block to start the blockchain!</p>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="text-left text-gray-600 text-sm border-b border-gray-200">
                    <th className="pb-4 font-medium">Block</th>
                    <th className="pb-4 font-medium">Hash</th>
                    <th className="pb-4 font-medium">Txs</th>
                    <th className="pb-4 font-medium">Miner</th>
                    <th className="pb-4 font-medium">Time</th>
                  </tr>
                </thead>
                <tbody>
                  {recentBlocks.map((block) => (
                    <tr key={block.hash} className="border-b border-gray-100">
                      <td className="py-4">
                        <span className="px-3 py-1 bg-purple-100 text-purple-600 rounded-lg font-semibold">
                          #{block.index}
                        </span>
                      </td>
                      <td className="py-4">
                        <code className="text-gray-500 text-sm">{block.hash?.substring(0, 16)}...</code>
                      </td>
                      <td className="py-4 text-gray-800">{block.transactionCount}</td>
                      <td className="py-4">
                        {block.minerWalletId === 'system' ? (
                          <span className="text-yellow-600">🌟 Genesis</span>
                        ) : (
                          <code className="text-gray-500 text-sm">{block.minerWalletId?.substring(0, 12)}...</code>
                        )}
                      </td>
                      <td className="py-4 text-gray-600 text-sm">
                        {new Date(block.timestamp).toLocaleString()}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default Mining;
