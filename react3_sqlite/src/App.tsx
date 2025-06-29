import React, { useState, useEffect } from 'react';
import CrudForm from './components/CrudForm';
import DataList from './components/DataList';
import { itemsApi} from './client/api/items';

function App() {
  const [data, setData] = useState([]);
  const [editingIndex, setEditingIndex] = useState(null);
  const [showForm, setShowForm] = useState(false);

  const fetchItems = async () => {
    try {
      //setLoading(true);
      const data = await itemsApi.getAll();
      console.log(data);
      setData(data);
    } catch (err) {
      //setError('アイテムの取得に失敗しました');
    } finally {
      //setLoading(false);
    }
  };

  useEffect(() => {
    fetchItems();
  }, []);

  const handleSubmit = async(formData) => {
    formData.is_public = 0;
    formData.food_orange = 0;
    formData.food_apple = 0;
    formData.food_banana = 0;
    formData.food_melon = 0;
    formData.food_grape = 0;
    console.log(formData);

    if (editingIndex !== null) {
      console.log("editingIndex=", editingIndex);
      const updatedData = [...data];
      const target = updatedData[editingIndex]
      console.log(target);
      await itemsApi.update(target.id, formData);
      //updatedData[editingIndex] = formData;
      //setData(updatedData);
    } else {
      await itemsApi.create(formData);
      //setData([...data, formData]);
    }
    await fetchItems();

    setShowForm(false);
    setEditingIndex(null);
  };

  const handleEdit = (index) => {
    setEditingIndex(index);
    setShowForm(true);
  };

  const handleDelete = async (index) => {
    if (window.confirm('本当に削除しますか？')) {
      //const updatedData = data.filter((_, i) => i !== index);
      const updatedData = data.filter((_, i) => i === index);
      if(updatedData.length > 0){
        const target = updatedData[0];
        //console.log(target);
        await itemsApi.delete(target.id);
        await fetchItems();
      }
    }
  };

  const handleCancel = () => {
    setShowForm(false);
    setEditingIndex(null);
  };

  const handleNewItem = () => {
    setEditingIndex(null);
    setShowForm(true);
  };

  return (
    <div className="min-h-screen bg-gray-100 py-8">
      <div className="container mx-auto px-4">
        <h1 className="text-3xl font-bold text-center mb-8 text-gray-800">
          CRUDアプリケーション
        </h1>
        
        {!showForm ? (
          <>
            <div className="text-center mb-8">
              <button
                onClick={handleNewItem}
                className="bg-green-600 text-white px-6 py-3 rounded-md hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 transition duration-200"
              >
                新規作成
              </button>
            </div>
            
            <DataList 
              data={data}
              onEdit={handleEdit}
              onDelete={handleDelete}
            />
          </>
        ) : (
          <CrudForm
            onSubmit={handleSubmit}
            initialData={editingIndex !== null ? data[editingIndex] : null}
            onCancel={handleCancel}
          />
        )}
      </div>
    </div>
  );
}

export default App;
