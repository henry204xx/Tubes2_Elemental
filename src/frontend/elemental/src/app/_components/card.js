"use client"

export default function Card({ title, image, tier, onClickedButt }) {
    if (tier == 0) return null;
    return (
        <div className="shadow-md rounded-lg p-4 w-50 bg-white hover:shadow-lg transition-shadow flex flex-col">
            <img src={image} alt={title} className="w-full h-30 object-contain mb-2" />
            <h2 className="text-xl font-semibold">{title}</h2>
            <p className="text-gray-500 mb-3">Tier: {tier}</p>
            <button 
                className="bg-green-500 hover:bg-green-600 text-white py-2 px-4 rounded-md transition-colors mt-auto"
                onClick={onClickedButt}
                suppressHydrationWarning
            >
                See My Receipt
            </button>
        </div>
    );
}