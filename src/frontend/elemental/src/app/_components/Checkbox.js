"use client";
import { useState } from "react";

export default function MultiValueToggle({ onChange }) {
  const [isMultivalue, setIsMultivalue] = useState(false);

  const handleChange = () => {
    const newValue = !isMultivalue;
    setIsMultivalue(newValue);
    onChange(newValue); // Kirim ke parent
  };

  return (
    <label className="flex items-center gap-2 cursor-pointer">
      <input
        type="checkbox"
        checked={isMultivalue}
        onChange={handleChange}
        className="accent-blue-500"
      />
      <span>Aktifkan Multivalue</span>
    </label>
  );
}
