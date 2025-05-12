"use client";

export default function MultiValueToggle({ onChange, checked }){

  const handleChange = () => {
    const newValue = !checked;
    onChange(newValue); 
  };

  return (
    <label className="flex items-center gap-2 cursor-pointer">
      <input
        type="checkbox"
        checked={checked}
        onChange={handleChange}
        className="accent-blue-500"
      />
      <span>Aktifkan Multivalue</span>
    </label>
  );
}
