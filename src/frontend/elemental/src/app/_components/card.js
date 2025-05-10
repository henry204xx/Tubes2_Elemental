"use client"

export default function Card({ title, image, tier }) {
    if (tier == 0) return null;
  return (
    <div className="shadow-md rounded-lg p-4 w-50 bg-white hover:shadow-lg transition-shadow">
      <img src={image} alt={title} className="w-full h-30 object-contain mb-2" />
      <h2 className="text-xl font-semibold">{title}</h2>
      <p className="text-gray-500">Tier: {tier}</p>
    </div>
  );
}