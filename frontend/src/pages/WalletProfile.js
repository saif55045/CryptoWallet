import React, { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import api from '../services/api';
import Navbar from '../components/Navbar';

const WalletProfile = () => {
  const { token } = useAuth();
  const [wallet, setWallet] = useState(null);
  const [balance, setBalance] = useState(0);
  const [utxos, setUtxos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [generating, setGenerating] = useState(false);
  const [showPrivateKey, setShowPrivateKey] = useState(false);
  const [privateKey, setPrivateKey] = useState('');
  const [copied, setCopied] = useState('');
  const [error, setError] = useState('');

  useEffect(() => {
    fetchWalletData();
  }, [token]);

  const fetchWalletData = async () => {
    try {
      const walletRes = await api.wallet.getMyWallet();
      setWallet(walletRes.data.wallet);

      const balanceRes = await api.utxo.getMyBalance();
      // Balance is returned as { balance: { walletId, balance, confirmedBalance, ... } }
      const balanceValue = balanceRes.data.balance?.balance || 0;
      setBalance(Number(balanceValue) || 0);

      const utxoRes = await api.utxo.getMyUTXOs();
      setUtxos(utxoRes.data.utxos || []);
    } catch (err) {
      if (err.response?.status !== 404) {
        setError('Failed to load wallet data');
      }
    } finally {
      setLoading(false);
    }
  };

  const handleGenerateWallet = async () => {
    setGenerating(true);
    setError('');
    try {
      const res = await api.wallet.generateWallet();
      setWallet(res.data.wallet);
      setPrivateKey(res.data.privateKey);
      setShowPrivateKey(true);
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to generate wallet');
    } finally {
      setGenerating(false);
    }
  };

  const handleExportKey = async () => {
    try {
      const res = await api.wallet.exportPrivateKey();
      setPrivateKey(res.data.privateKey);
      setShowPrivateKey(true);
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to export private key');
    }
  };

  const copyToClipboard = (text, type) => {
    navigator.clipboard.writeText(text);
    setCopied(type);
    setTimeout(() => setCopied(''), 2000);
  };

  if (loading) {
    return (
      <div className="min-h-screen">
        <Navbar />
        <div className="flex items-center justify-center h-96">
          <div className="text-center">
            <div className="w-16 h-16 border-4 border-indigo-500 border-t-transparent rounded-full animate-spin mx-auto"></div>
            <p className="text-gray-600 mt-4">Loading wallet...</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen">
      <Navbar />
      
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-800 mb-2">💼 Wallet</h1>
          <p className="text-gray-600">Manage your crypto wallet and keys</p>
        </div>

        {error && (
          <div className="mb-6 bg-red-500/10 border border-red-500/50 text-red-400 px-4 py-3 rounded-xl">
            {error}
          </div>
        )}

        {!wallet ? (
          /* No Wallet State */
          <div className="glass-card rounded-2xl p-12 text-center">
            <div className="w-24 h-24 bg-gradient-to-br from-indigo-500/20 to-purple-500/20 rounded-3xl flex items-center justify-center mx-auto mb-6">
              <span className="text-5xl">💼</span>
            </div>
            <h2 className="text-2xl font-bold text-gray-800 mb-3">Create Your Wallet</h2>
            <p className="text-gray-600 mb-8 max-w-md mx-auto">
              Generate a secure ECDSA key pair to start sending and receiving crypto
            </p>
            <button
              onClick={handleGenerateWallet}
              disabled={generating}
              className="btn-primary inline-flex items-center gap-3 text-lg disabled:opacity-50"
            >
              {generating ? (
                <>
                  <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                  </svg>
                  Generating...
                </>
              ) : (
                <>
                  <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                  </svg>
                  Generate Wallet
                </>
              )}
            </button>
          </div>
        ) : (
          <div className="space-y-6">
            {/* Balance Card */}
            <div className="glass-card rounded-2xl p-8 bg-gradient-to-br from-indigo-500/10 to-purple-500/10">
              <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-6">
                <div>
                  <p className="text-gray-600 text-sm font-medium mb-1">Total Balance</p>
                  <p className="text-5xl font-bold text-gray-800">{balance.toFixed(4)}</p>
                  <p className="text-gray-600 mt-1">coins</p>
                </div>
                <div className="flex items-center gap-2">
                  <div className="w-3 h-3 bg-green-500 rounded-full animate-pulse"></div>
                  <span className="text-green-400 font-medium">Wallet Active</span>
                </div>
              </div>
            </div>

            {/* Wallet Address */}
            <div className="glass-card rounded-2xl p-6">
              <h3 className="text-lg font-semibold text-gray-800 mb-4 flex items-center gap-2">
                <svg className="w-5 h-5 text-indigo-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
                </svg>
                Wallet Address
              </h3>
              <div className="bg-gray-50 rounded-xl p-4 flex items-center justify-between gap-4">
                <code className="text-gray-700 text-sm break-all flex-1">{wallet.walletId}</code>
                <button
                  onClick={() => copyToClipboard(wallet.walletId, 'address')}
                  className="p-2 rounded-lg bg-white/10 hover:bg-white/20 transition-colors flex-shrink-0"
                >
                  {copied === 'address' ? (
                    <svg className="w-5 h-5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                  ) : (
                    <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                    </svg>
                  )}
                </button>
              </div>
            </div>

            {/* Public Key */}
            <div className="glass-card rounded-2xl p-6">
              <h3 className="text-lg font-semibold text-gray-800 mb-4 flex items-center gap-2">
                <svg className="w-5 h-5 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                </svg>
                Public Key
              </h3>
              <div className="bg-gray-50 rounded-xl p-4 flex items-center justify-between gap-4">
                <code className="text-gray-700 text-xs break-all flex-1">{wallet.publicKey}</code>
                <button
                  onClick={() => copyToClipboard(wallet.publicKey, 'public')}
                  className="p-2 rounded-lg bg-white/10 hover:bg-white/20 transition-colors flex-shrink-0"
                >
                  {copied === 'public' ? (
                    <svg className="w-5 h-5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                  ) : (
                    <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                    </svg>
                  )}
                </button>
              </div>
            </div>

            {/* Private Key Section */}
            <div className="glass-card rounded-2xl p-6 border border-red-500/20">
              <h3 className="text-lg font-semibold text-gray-800 mb-4 flex items-center gap-2">
                <svg className="w-5 h-5 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
                Private Key
                <span className="text-xs bg-red-500/20 text-red-400 px-2 py-0.5 rounded-full ml-2">Sensitive</span>
              </h3>
              
              {showPrivateKey && privateKey ? (
                <div className="space-y-4">
                  <div className="bg-red-500/10 border border-red-500/30 rounded-xl p-4">
                    <div className="flex items-start gap-2 text-red-400 text-sm mb-3">
                      <svg className="w-5 h-5 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
                        <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                      </svg>
                      Never share your private key! Anyone with this key can access your funds.
                    </div>
                    <code className="text-gray-700 text-xs break-all block">{privateKey}</code>
                  </div>
                  <div className="flex gap-3">
                    <button
                      onClick={() => copyToClipboard(privateKey, 'private')}
                      className="flex-1 px-4 py-2 rounded-xl bg-gray-100 hover:bg-gray-200 text-gray-800 text-sm font-medium transition-colors flex items-center justify-center gap-2"
                    >
                      {copied === 'private' ? '✓ Copied!' : 'Copy Key'}
                    </button>
                    <button
                      onClick={() => setShowPrivateKey(false)}
                      className="px-4 py-2 rounded-xl bg-red-500/20 hover:bg-red-500/30 text-red-400 text-sm font-medium transition-colors"
                    >
                      Hide
                    </button>
                  </div>
                </div>
              ) : (
                <button
                  onClick={handleExportKey}
                  className="w-full px-4 py-3 rounded-xl bg-gray-100 hover:bg-gray-200 text-gray-700 font-medium transition-colors flex items-center justify-center gap-2"
                >
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                  </svg>
                  Reveal Private Key
                </button>
              )}
            </div>

            {/* UTXOs */}
            <div className="glass-card rounded-2xl p-6">
              <h3 className="text-lg font-semibold text-gray-800 mb-4 flex items-center gap-2">
                <svg className="w-5 h-5 text-purple-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                </svg>
                Unspent Outputs (UTXOs)
                <span className="text-xs bg-purple-100 text-purple-600 px-2 py-0.5 rounded-full ml-2">{utxos.length}</span>
              </h3>
              
              {utxos.length > 0 ? (
                <div className="space-y-2 max-h-64 overflow-y-auto">
                  {utxos.map((utxo, index) => (
                    <div key={index} className="bg-gray-50 rounded-xl p-4 flex items-center justify-between">
                      <div>
                        <p className="text-gray-600 text-xs">TX: {utxo.transactionId?.substring(0, 16)}...</p>
                        <p className="text-gray-500 text-xs">Output #{utxo.outputIndex}</p>
                      </div>
                      <div className="text-right">
                        <p className="text-green-600 font-semibold">{utxo.amount?.toFixed(4)}</p>
                        <p className={`text-xs ${utxo.isConfirmed ? 'text-green-600' : 'text-yellow-600'}`}>
                          {utxo.isConfirmed ? 'Confirmed' : 'Pending'}
                        </p>
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="text-center py-6 text-gray-600">
                  <p>No UTXOs yet. Mine some blocks to get coins!</p>
                </div>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default WalletProfile;
