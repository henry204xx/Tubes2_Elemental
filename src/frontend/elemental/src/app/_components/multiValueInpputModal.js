"use client"

export default function MultiValueInputModal({
  isOpen,
  onClose,
  setInputValue,
  onSubmit,
  inputError,
}) {
  if (!isOpen) return null;

  return (
    <div 
      className="fixed inset-0 flex items-center justify-center z-50"
      style={{ backgroundColor: 'rgba(209, 213, 219, 0.7)' }} 
    >
      <div className="bg-white p-6 rounded-lg shadow-lg max-w-sm w-full">
        <h2 className="text-xl font-bold text-black mb-4">Masukkan Nilai</h2>
        <input
          type="number"
          onChange={(e) => setInputValue(parseInt(e.target.value))}
          className="border p-2 w-full text-black mb-4"
          placeholder="Masukkan nilai lebih dari 0"
        />
        {inputError && <p className="text-red-500 text-sm mb-4">{inputError}</p>}
        <div className="flex justify-between">
          <button
            onClick={onClose}
            className="px-4 py-2 bg-gray-300 text-black rounded-md"
          >
            Cancel
          </button>
          <button
            onClick={onSubmit}
            className="px-4 py-2 bg-blue-500 text-white rounded-md"
          >
            OK
          </button>
        </div>
      </div>
    </div>
  );
}
