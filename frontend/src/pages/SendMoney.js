import React, { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import api from '../services/api';
import Navbar from '../components/Navbar';

const SendMoney = () => {
  const { token } = useAuth();
  const [wallet, setWallet] = useState(null);
  const [balance, setBalance] = useState(0);
  const [beneficiaries, setBeneficiaries] = useState([]);
  const [transactions, setTransactions] = useState([]);
  const [formData, setFormData] = useState({
    recipientWalletId: '',
    amount: '',
    message: '',
  });
  const [loading, setLoading] = useState(true);
  const [sending, setSending] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [activeTab, setActiveTab] = useState('send');

  useEffect(() => {
    fetchData();
  }, [token]);

  const fetchData = async () => {
    try {
      const [walletRes, balanceRes, beneficiariesRes, txRes] = await Promise.all([
        api.wallet.getMyWallet(),
        api.utxo.getMyBalance(),
        api.wallet.getBeneficiaries(),
        api.transaction.getMyTransactions(),
      ]);
      
      setWallet(walletRes.data.wallet);
      // Balance is returned as { balance: { walletId, balance, confirmedBalance, ... } }
      const balanceValue = balanceRes.data.balance?.balance || 0;
      setBalance(Number(balanceValue) || 0);
      setBeneficiaries(beneficiariesRes.data.beneficiaries || []);
      setTransactions(txRes.data.transactions || []);
    } catch (err) {
      console.log('Loading data...');
      setBalance(0);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setSending(true);
    setError('');
    setSuccess('');

    try {
      await api.transaction.send({
        recipientWalletId: formData.recipientWalletId,
        amount: parseFloat(formData.amount),
        message: formData.message,
      });
      
      setSuccess('Transaction sent successfully! It will be confirmed when mined.');
      setFormData({ recipientWallet: '', amount: '', message: '' });
      fetchData();
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to send transaction');
    } finally {
      setSending(false);
    }
  };

  const selectBeneficiary = (walletId) => {
    setFormData({ ...formData, recipientWalletId: walletId });
    setActiveTab('send');
  };

  if (loading) {
    return (
      <div className="min-h-screen">
        <Navbar />
        <div className="flex items-center justify-center h-96">
          <div className="text-center">
            <div className="w-16 h-16 border-4 border-indigo-500 border-t-transparent rounded-full animate-spin mx-auto"></div>
            <p className="text-gray-600 mt-4">Loading...</p>
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
          <h1 className="text-3xl font-bold text-gray-800 mb-2">📤 Send Money</h1>
          <p className="text-gray-600">Transfer coins securely using digital signatures</p>
        </div>

        {/* Balance Banner */}
        <div className="glass-card rounded-2xl p-6 mb-8 bg-gradient-to-r from-indigo-500/10 to-purple-500/10">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-gray-600 text-sm">Available Balance</p>
              <p className="text-3xl font-bold text-gray-800">{balance.toFixed(4)} <span className="text-lg text-gray-500">coins</span></p>
            </div>
            <div className="w-14 h-14 rounded-2xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center">
              <span className="text-2xl">💰</span>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Send Form */}
          <div className="lg:col-span-2">
            {/* Tabs */}
            <div className="flex gap-2 mb-6">
              <button
                onClick={() => setActiveTab('send')}
                className={`px-4 py-2 rounded-xl font-medium transition-all ${
                  activeTab === 'send'
                    ? 'bg-indigo-500 text-white'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                }`}
              >
                Send
              </button>
              <button
                onClick={() => setActiveTab('contacts')}
                className={`px-4 py-2 rounded-xl font-medium transition-all ${
                  activeTab === 'contacts'
                    ? 'bg-indigo-500 text-white'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                }`}
              >
                Contacts ({beneficiaries.length})
              </button>
              <button
                onClick={() => setActiveTab('history')}
                className={`px-4 py-2 rounded-xl font-medium transition-all ${
                  activeTab === 'history'
                    ? 'bg-indigo-500 text-white'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                }`}
              >
                History
              </button>
            </div>

            {activeTab === 'send' && (
              <div className="glass-card rounded-2xl p-6">
                <h2 className="text-xl font-bold text-gray-800 mb-6">Send Transaction</h2>

                {error && (
                  <div className="mb-6 bg-red-500/10 border border-red-500/50 text-red-400 px-4 py-3 rounded-xl flex items-center gap-2">
                    <svg className="w-5 h-5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                    </svg>
                    {error}
                  </div>
                )}

                {success && (
                  <div className="mb-6 bg-green-500/10 border border-green-500/50 text-green-400 px-4 py-3 rounded-xl flex items-center gap-2">
                    <svg className="w-5 h-5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                    </svg>
                    {success}
                  </div>
                )}

                <form onSubmit={handleSubmit} className="space-y-6">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Recipient Wallet Address
                    </label>
                    <input
                      type="text"
                      value={formData.recipientWalletId}
                      onChange={(e) => setFormData({ ...formData, recipientWalletId: e.target.value })}
                      className="input-modern"
                      placeholder="Enter wallet address"
                      required
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Amount
                    </label>
                    <div className="relative">
                      <input
                        type="number"
                        step="0.0001"
                        min="0.0001"
                        max={balance}
                        value={formData.amount}
                        onChange={(e) => setFormData({ ...formData, amount: e.target.value })}
                        className="input-modern pr-20"
                        placeholder="0.00"
                        required
                      />
                      <button
                        type="button"
                        onClick={() => setFormData({ ...formData, amount: balance.toString() })}
                        className="absolute right-2 top-1/2 -translate-y-1/2 px-3 py-1 bg-indigo-500/20 text-indigo-400 text-sm rounded-lg hover:bg-indigo-500/30 transition-colors"
                      >
                        MAX
                      </button>
                    </div>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Message (Optional)
                    </label>
                    <input
                      type="text"
                      value={formData.message}
                      onChange={(e) => setFormData({ ...formData, message: e.target.value })}
                      className="input-modern"
                      placeholder="Add a note..."
                    />
                  </div>

                  <button
                    type="submit"
                    disabled={sending || !formData.recipientWalletId || !formData.amount}
                    className="w-full btn-primary flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {sending ? (
                      <>
                        <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
                          <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                          <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                        </svg>
                        Signing & Sending...
                      </>
                    ) : (
                      <>
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
                        </svg>
                        Send Transaction
                      </>
                    )}
                  </button>
                </form>
              </div>
            )}

            {activeTab === 'contacts' && (
              <div className="glass-card rounded-2xl p-6">
                <h2 className="text-xl font-bold text-gray-800 mb-6">Saved Contacts</h2>
                
                {beneficiaries.length > 0 ? (
                  <div className="space-y-3">
                    {beneficiaries.map((b) => (
                      <div
                        key={b._id}
                        className="bg-gray-50 rounded-xl p-4 flex items-center justify-between hover:bg-gray-100 transition-colors cursor-pointer"
                        onClick={() => selectBeneficiary(b.walletId)}
                      >
                        <div className="flex items-center gap-3">
                          <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-semibold">
                            {b.name?.charAt(0).toUpperCase()}
                          </div>
                          <div>
                            <p className="text-gray-800 font-medium">{b.name}</p>
                            <p className="text-gray-500 text-sm font-mono">{b.walletId?.substring(0, 20)}...</p>
                          </div>
                        </div>
                        <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                        </svg>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-8 text-gray-600">
                    <p>No saved contacts yet.</p>
                  </div>
                )}
              </div>
            )}

            {activeTab === 'history' && (
              <div className="glass-card rounded-2xl p-6">
                <h2 className="text-xl font-bold text-gray-800 mb-6">Transaction History</h2>
                
                {transactions.length > 0 ? (
                  <div className="space-y-3">
                    {transactions.slice(0, 10).map((tx) => (
                      <div key={tx.transactionId} className="bg-gray-50 rounded-xl p-4">
                        <div className="flex items-center justify-between mb-2">
                          <div className="flex items-center gap-3">
                            <div className={`w-10 h-10 rounded-xl flex items-center justify-center ${
                              tx.direction === 'sent'
                                ? 'bg-red-500/20 text-red-400'
                                : 'bg-green-500/20 text-green-400'
                            }`}>
                              {tx.direction === 'sent' ? '↑' : '↓'}
                            </div>
                            <div>
                              <p className="text-gray-800 font-medium">
                                {tx.direction === 'sent' ? 'Sent' : 'Received'}
                              </p>
                              <p className="text-gray-500 text-xs">
                                {new Date(tx.timestamp).toLocaleString()}
                              </p>
                            </div>
                          </div>
                          <div className="text-right">
                            <p className={`font-semibold ${
                              tx.direction === 'sent' ? 'text-red-400' : 'text-green-400'
                            }`}>
                              {tx.direction === 'sent' ? '-' : '+'}{tx.amount?.toFixed(4)}
                            </p>
                            <span className={`text-xs px-2 py-0.5 rounded-full ${
                              tx.status === 'confirmed'
                                ? 'bg-green-500/20 text-green-400'
                                : 'bg-yellow-500/20 text-yellow-400'
                            }`}>
                              {tx.status}
                            </span>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-8 text-gray-600">
                    <p>No transactions yet.</p>
                  </div>
                )}
              </div>
            )}
          </div>

          {/* Info Sidebar */}
          <div className="space-y-6">
            <div className="glass-card rounded-2xl p-6">
              <h3 className="text-lg font-semibold text-gray-800 mb-4">How it works</h3>
              <div className="space-y-4">
                <div className="flex gap-3">
                  <div className="w-8 h-8 rounded-lg bg-indigo-500/20 flex items-center justify-center text-indigo-600 font-semibold flex-shrink-0">1</div>
                  <div>
                    <p className="text-gray-800 font-medium">Enter Details</p>
                    <p className="text-gray-600 text-sm">Recipient address and amount</p>
                  </div>
                </div>
                <div className="flex gap-3">
                  <div className="w-8 h-8 rounded-lg bg-indigo-500/20 flex items-center justify-center text-indigo-600 font-semibold flex-shrink-0">2</div>
                  <div>
                    <p className="text-gray-800 font-medium">Sign Transaction</p>
                    <p className="text-gray-600 text-sm">ECDSA signature created</p>
                  </div>
                </div>
                <div className="flex gap-3">
                  <div className="w-8 h-8 rounded-lg bg-indigo-500/20 flex items-center justify-center text-indigo-600 font-semibold flex-shrink-0">3</div>
                  <div>
                    <p className="text-gray-800 font-medium">Broadcast</p>
                    <p className="text-gray-600 text-sm">Sent to the network</p>
                  </div>
                </div>
                <div className="flex gap-3">
                  <div className="w-8 h-8 rounded-lg bg-green-500/20 flex items-center justify-center text-green-600 font-semibold flex-shrink-0">4</div>
                  <div>
                    <p className="text-gray-800 font-medium">Confirmation</p>
                    <p className="text-gray-600 text-sm">Mined into a block</p>
                  </div>
                </div>
              </div>
            </div>

            <div className="glass-card rounded-2xl p-6">
              <h3 className="text-lg font-semibold text-gray-800 mb-4">Your Wallet</h3>
              <div className="bg-gray-50 rounded-xl p-3">
                <p className="text-gray-600 text-xs mb-1">Address</p>
                <p className="text-gray-800 font-mono text-xs break-all">{wallet?.walletId}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SendMoney;
