import React, { useState, useEffect } from 'react';
//import './App.css';

function App() {
  const [todos, setTodos] = useState([]);
  const [showModal, setShowModal] = useState(false);
  const [isEdit, setIsEdit] = useState(false);
  const [currentTodo, setCurrentTodo] = useState(null);

  useEffect(() => {
    fetch('/api/list')
      .then(res => res.json())
      .then(data => setTodos(data));
  }, []);

  const handleOpenModal = (todo = null) => {
    if (todo) {
      setIsEdit(true);
      setCurrentTodo(todo);
    } else {
      setIsEdit(false);
      setCurrentTodo({
        title: '',
        content: '',
        completed: false,
        content_type: '',
        is_public: false,
        food_orange: false,
        food_apple: false,
        food_banana: false,
        food_melon: false,
        food_grape: false,
        pub_date1: '',
        pub_date2: '',
        pub_date3: '',
        pub_date4: '',
        pub_date5: '',
        pub_date6: '',
        qty1: '',
        qty2: '',
        qty3: '',
        qty4: '',
        qty5: '',
        qty6: '',
      });
    }
    setShowModal(true);
  };

  const handleCloseModal = () => {
    setShowModal(false);
  };

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setCurrentTodo(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value,
    }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    const url = isEdit ? `/api/update` : '/api/create';
    const method = 'POST';

    fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(currentTodo),
    })
      .then(res => res.json())
      .then(data => {
        if (isEdit) {
          setTodos(todos.map(todo => (todo.id === data.id ? data : todo)));
        } else {
          setTodos([data, ...todos]);
        }
        handleCloseModal();
      });
  };

  const handleDelete = (id) => {
    fetch('/api/delete', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ id }),
    }).then(() => {
      setTodos(todos.filter(todo => todo.id !== id));
    });
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Todo App</h1>
      <button
        onClick={() => handleOpenModal()}
        className="bg-blue-500 text-white px-4 py-2 rounded mb-4"
      >
        Add Todo
      </button>

      {showModal && (
        <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full">
          <div className="relative top-20 mx-auto p-5 border w-1/2 shadow-lg rounded-md bg-white">
            <h3 className="text-lg font-bold">{isEdit ? 'Edit Todo' : 'Add Todo'}</h3>
            <form onSubmit={handleSubmit}>
              <div className="grid grid-cols-2 gap-4">
                <input type="text" name="title" value={currentTodo.title} onChange={handleChange} placeholder="Title" className="border p-2" />
                <input type="text" name="content" value={currentTodo.content} onChange={handleChange} placeholder="Content" className="border p-2" />
                <label><input type="checkbox" name="completed" checked={currentTodo.completed} onChange={handleChange} /> Completed</label>
                <input type="text" name="content_type" value={currentTodo.content_type} onChange={handleChange} placeholder="Content Type" className="border p-2" />
                <div>
                  <label><input type="radio" name="is_public" value="true" checked={currentTodo.is_public === true} onChange={() => setCurrentTodo(prev => ({ ...prev, is_public: true }))} /> Public</label>
                  <label><input type="radio" name="is_public" value="false" checked={currentTodo.is_public === false} onChange={() => setCurrentTodo(prev => ({ ...prev, is_public: false }))} /> Private</label>
                </div>
                <div>
                  <label><input type="checkbox" name="food_orange" checked={currentTodo.food_orange} onChange={handleChange} /> Orange</label>
                  <label><input type="checkbox" name="food_apple" checked={currentTodo.food_apple} onChange={handleChange} /> Apple</label>
                  <label><input type="checkbox" name="food_banana" checked={currentTodo.food_banana} onChange={handleChange} /> Banana</label>
                  <label><input type="checkbox" name="food_melon" checked={currentTodo.food_melon} onChange={handleChange} /> Melon</label>
                  <label><input type="checkbox" name="food_grape" checked={currentTodo.food_grape} onChange={handleChange} /> Grape</label>
                </div>
                <input type="date" name="pub_date1" value={currentTodo.pub_date1} onChange={handleChange} className="border p-2" />
                <input type="date" name="pub_date2" value={currentTodo.pub_date2} onChange={handleChange} className="border p-2" />
                <input type="date" name="pub_date3" value={currentTodo.pub_date3} onChange={handleChange} className="border p-2" />
                <input type="date" name="pub_date4" value={currentTodo.pub_date4} onChange={handleChange} className="border p-2" />
                <input type="date" name="pub_date5" value={currentTodo.pub_date5} onChange={handleChange} className="border p-2" />
                <input type="date" name="pub_date6" value={currentTodo.pub_date6} onChange={handleChange} className="border p-2" />
                <input type="text" name="qty1" value={currentTodo.qty1} onChange={handleChange} placeholder="Qty 1" className="border p-2" />
                <input type="text" name="qty2" value={currentTodo.qty2} onChange={handleChange} placeholder="Qty 2" className="border p-2" />
                <input type="text" name="qty3" value={currentTodo.qty3} onChange={handleChange} placeholder="Qty 3" className="border p-2" />
                <input type="text" name="qty4" value={currentTodo.qty4} onChange={handleChange} placeholder="Qty 4" className="border p-2" />
                <input type="text" name="qty5" value={currentTodo.qty5} onChange={handleChange} placeholder="Qty 5" className="border p-2" />
                <input type="text" name="qty6" value={currentTodo.qty6} onChange={handleChange} placeholder="Qty 6" className="border p-2" />
              </div>
              <div className="flex justify-end mt-4">
                <button type="button" onClick={handleCloseModal} className="bg-gray-500 text-white px-4 py-2 rounded mr-2">Cancel</button>
                <button type="submit" className="bg-blue-500 text-white px-4 py-2 rounded">Submit</button>
              </div>
            </form>
          </div>
        </div>
      )}

      <table className="table-auto w-full">
        <thead>
          <tr>
            <th className="px-4 py-2">ID</th>
            <th className="px-4 py-2">Title</th>
            <th className="px-4 py-2">Content</th>
            <th className="px-4 py-2">Completed</th>
            <th className="px-4 py-2">Actions</th>
          </tr>
        </thead>
        <tbody>
          {todos.map(todo => (
            <tr key={todo.id}>
              <td className="border px-4 py-2">{todo.id}</td>
              <td className="border px-4 py-2">{todo.title}</td>
              <td className="border px-4 py-2">{todo.content}</td>
              <td className="border px-4 py-2">{todo.completed ? 'Yes' : 'No'}</td>
              <td className="border px-4 py-2">
                <button
                  onClick={() => handleOpenModal(todo)}
                  className="bg-green-500 text-white px-2 py-1 rounded mr-2"
                >
                  Edit
                </button>
                <button
                  onClick={() => handleDelete(todo.id)}
                  className="bg-red-500 text-white px-2 py-1 rounded"
                >
                  Delete
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default App;