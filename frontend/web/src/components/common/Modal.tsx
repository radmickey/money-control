import React from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { X } from 'lucide-react';

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  children: React.ReactNode;
  footer?: React.ReactNode;
}

const Modal: React.FC<ModalProps> = ({ isOpen, onClose, title, children, footer }) => {
  if (!isOpen) return null;

  return (
    <AnimatePresence>
      <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <motion.div
          initial={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          exit={{ opacity: 0, scale: 0.95 }}
          className="glass rounded-2xl p-6 w-full max-w-md"
        >
          <div className="flex items-center justify-between mb-6">
            <h3 className="text-xl font-semibold">{title}</h3>
            <button
              onClick={onClose}
              className="p-2 text-midnight-400 hover:text-white transition-colors"
            >
              <X className="w-5 h-5" />
            </button>
          </div>

          <div className="space-y-4">{children}</div>

          {footer && <div className="flex gap-3 mt-6">{footer}</div>}
        </motion.div>
      </div>
    </AnimatePresence>
  );
};

export default Modal;

// Reusable form field components
export const FormField: React.FC<{
  label: string;
  children: React.ReactNode;
}> = ({ label, children }) => (
  <div>
    <label className="block text-sm text-midnight-400 mb-2">{label}</label>
    {children}
  </div>
);

// Reusable button components
export const CancelButton: React.FC<{ onClick: () => void }> = ({ onClick }) => (
  <button
    onClick={onClick}
    className="flex-1 px-4 py-3 rounded-xl border border-midnight-700 text-midnight-300 hover:bg-midnight-800/50 transition-colors"
  >
    Cancel
  </button>
);

export const SubmitButton: React.FC<{
  onClick: () => void;
  disabled?: boolean;
  children: React.ReactNode;
}> = ({ onClick, disabled, children }) => (
  <button
    onClick={onClick}
    disabled={disabled}
    className="flex-1 btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
  >
    {children}
  </button>
);

// Currency selector
export const CurrencySelect: React.FC<{
  value: string;
  onChange: (value: string) => void;
}> = ({ value, onChange }) => (
  <select
    value={value}
    onChange={(e) => onChange(e.target.value)}
    className="input-field"
  >
    <option value="USD">USD - US Dollar</option>
    <option value="EUR">EUR - Euro</option>
    <option value="GBP">GBP - British Pound</option>
    <option value="RUB">RUB - Russian Ruble</option>
    <option value="JPY">JPY - Japanese Yen</option>
    <option value="CHF">CHF - Swiss Franc</option>
    <option value="CAD">CAD - Canadian Dollar</option>
    <option value="AUD">AUD - Australian Dollar</option>
  </select>
);

