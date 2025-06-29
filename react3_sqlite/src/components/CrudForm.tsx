import React, { useState } from 'react';

const CrudForm = ({ onSubmit, initialData = null, onCancel }) => {
  const [formData, setFormData] = useState({
    title: initialData?.title || '',
    content: initialData?.content || '',
    content_type: initialData?.content_type || '',
    is_public: initialData?.is_public || 'public',
    food_orange: initialData?.food_orange || false,
    food_apple: initialData?.food_apple || false,
    food_banana: initialData?.food_banana || false,
    food_melon: initialData?.food_melon || false,
    food_grape: initialData?.food_grape || false,
  });

  const [errors, setErrors] = useState({});

  const handleInputChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
  };

  const validateForm = () => {
    const newErrors = {};
    
    if (!formData.title.trim()) {
      newErrors.title = 'タイトルは必須です';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    if (validateForm()) {
      onSubmit(formData);
    }
  };

  return (
    <div className="max-w-2xl mx-auto p-6 bg-white rounded-lg shadow-md">
      <h2 className="text-2xl font-bold mb-6 text-gray-800">
        {initialData ? '編集' : '新規作成'}
      </h2>
      
      <form onSubmit={handleSubmit} className="space-y-6">
        <div>
          <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-1">
            タイトル *
          </label>
          <input
            type="text"
            id="title"
            name="title"
            value={formData.title}
            onChange={handleInputChange}
            className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
              errors.title ? 'border-red-500' : 'border-gray-300'
            }`}
            placeholder="タイトルを入力してください"
          />
          {errors.title && (
            <p className="text-red-500 text-sm mt-1">{errors.title}</p>
          )}
        </div>

        <div>
          <label htmlFor="content" className="block text-sm font-medium text-gray-700 mb-1">
            コンテンツ
          </label>
          <input
            type="text"
            id="content"
            name="content"
            value={formData.content}
            onChange={handleInputChange}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="コンテンツを入力してください"
          />
        </div>

        <div>
          <label htmlFor="content_type" className="block text-sm font-medium text-gray-700 mb-1">
            コンテンツタイプ
          </label>
          <input
            type="text"
            id="content_type"
            name="content_type"
            value={formData.content_type}
            onChange={handleInputChange}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="コンテンツタイプを入力してください"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-3">
            公開設定
          </label>
          <div className="space-y-2">
            <label className="flex items-center">
              <input
                type="radio"
                name="is_public"
                value="public"
                checked={formData.is_public === 'public'}
                onChange={handleInputChange}
                className="mr-2 text-blue-600 focus:ring-blue-500"
              />
              <span className="text-sm text-gray-700">公開</span>
            </label>
            <label className="flex items-center">
              <input
                type="radio"
                name="is_public"
                value="private"
                checked={formData.is_public === 'private'}
                onChange={handleInputChange}
                className="mr-2 text-blue-600 focus:ring-blue-500"
              />
              <span className="text-sm text-gray-700">非公開</span>
            </label>
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-3">
            好きなフルーツ
          </label>
          <div className="grid grid-cols-2 gap-3">
            {[
              { name: 'food_orange', label: 'オレンジ' },
              { name: 'food_apple', label: 'りんご' },
              { name: 'food_banana', label: 'バナナ' },
              { name: 'food_melon', label: 'メロン' },
              { name: 'food_grape', label: 'ぶどう' }
            ].map(fruit => (
              <label key={fruit.name} className="flex items-center">
                <input
                  type="checkbox"
                  name={fruit.name}
                  checked={formData[fruit.name]}
                  onChange={handleInputChange}
                  className="mr-2 text-blue-600 focus:ring-blue-500"
                />
                <span className="text-sm text-gray-700">{fruit.label}</span>
              </label>
            ))}
          </div>
        </div>

        <div className="flex gap-4 pt-4">
          <button
            type="submit"
            className="flex-1 bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition duration-200"
          >
            {initialData ? '更新' : '作成'}
          </button>
          
          {onCancel && (
            <button
              type="button"
              onClick={onCancel}
              className="flex-1 bg-gray-500 text-white py-2 px-4 rounded-md hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 transition duration-200"
            >
              キャンセル
            </button>
          )}
        </div>
      </form>
    </div>
  );
};

export default CrudForm;