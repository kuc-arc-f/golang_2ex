import React, { useState, useEffect } from 'react';
import {
  useReactTable,
  getCoreRowModel,
  getPaginationRowModel,
  flexRender,
  ColumnDef,
} from "@tanstack/react-table";

import { Item, NewItem } from './types/Item';
import { itemsApi } from './ContetData/api';
import dataUtil from './ContetData/dataUtil';
import ItemDialog from './ContetData/ItemDialog';
import Head from '../components/Head';
let contentId: number = 0;

function App() {
  const [items, setItems] = useState<Item[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [dialogMode, setDialogMode] = useState<'create' | 'edit'>('create');
  const [editingItem, setEditingItem] = useState<undefined>();
  const columns: ColumnDef<Person>[] = [
    { header: "ID", accessorKey: "id" },
    { header: "DataJson", accessorKey: "data_list" },
    {
      header: "操作",
      id: "display",
      cell: ({ row }) => {
        //console.log(row.original);
        return (
          <button
            onClick={() => handleEdit(row.original)}
            className="text-indigo-600 hover:text-indigo-900 mr-4"
          >
            [ Show ]
          </button>        
        )
      },
    },  
  ];

  const table = useReactTable({
    data: items,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(), // ページネーション追加
  });

  // アイテム一覧を取得
  const fetchItems = async () => {
    try {
      setLoading(true);
      const data = await itemsApi.getAll(contentId);
      console.log(data);
      setItems(data);
    } catch (err) {
      setError('アイテムの取得に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    const searchParams = new URLSearchParams(window.location.search);
    contentId = searchParams.get('content') || "";
    console.log("contentId=", contentId);
    fetchItems();
  }, []);

  // 新規作成ダイアログを開く
  const handleCreate = () => {
    setDialogMode('create');
    setEditingItem(undefined);
    setDialogOpen(true);
  };

  // 編集ダイアログを開く
  const handleEdit = (item: Item) => {
    setDialogMode('edit');
    setEditingItem(item);
    setDialogOpen(true);
  };

  // アイテム保存
  const handleSave = async (itemData: NewItem) => {
    try {
      console.log(itemData);
      if (dialogMode === 'create') {
        await itemsApi.create(itemData);
      } else if (editingItem) {
        await itemsApi.update(editingItem.id, itemData);
      }
      await fetchItems();
      setError(null);
    } catch (err) {
      setError('保存に失敗しました');
    }
  };

  // アイテム削除
  const handleDelete = async (id: number) => {
    if (!confirm('本当に削除しますか？')) return;
    
    try {
      await itemsApi.delete(id);
      await fetchItems();
      setError(null);
    } catch (err) {
      setError('削除に失敗しました');
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="text-xl">読み込み中...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100 pt-1 pb-8">
      <Head />
      <div className="max-w-6xl mx-auto px-4">
        <a href="/">
          <button
            className="bg-white text-black my-2 px-4 py-2 rounded-md hover:bg-blue-700"
          >
          Back
          </button>
        </a>
        <div className="bg-white rounded-lg shadow pb-8">
          <div className="px-6 py-4 border-b border-gray-200 flex justify-between items-center">
            <h1 className="text-2xl font-bold text-gray-900">Data</h1>
          </div>

          <div className="p-2">
            {items.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                アイテムがありません
              </div>
            ) : ""}
          </div>
          <table className="border border-gray-300 w-full">
            <thead>
              {table.getHeaderGroups().map((headerGroup) => (
                <tr key={headerGroup.id} className="bg-gray-100">
                  {headerGroup.headers.map((header) => (

                    <th key={header.id} className="border p-2 text-left">
                      {header.isPlaceholder
                        ? null
                        : flexRender(header.column.columnDef.header, header.getContext())}
                      {header.id === "name" ? (<span>
                        <button onClick={() => { sortStart("name", sortName); }}>
                          <span className="ms-2 text-green-600">sort</span>
                        </button>
                      </span>) : ""}
                      {header.id === "age" ? (<span>
                        <button onClick={() => { sortStart("age", sortAge); }}>
                          <span className="ms-2 text-green-600">sort</span>
                        </button>
                      </span>) : ""}                      
                      {header.id === "weight" ? (<span>
                        <button onClick={() => { sortStart("weight", sortWeight); }}>
                          <span className="ms-2 text-green-600">sort</span>
                        </button>
                      </span>) : ""}                      

                    </th>
                  ))}
                </tr>
              ))}
            </thead>
            <tbody>
              {table.getRowModel().rows.map((row) => (
                <tr key={row.id}>
                  {row.getVisibleCells().map((cell) => (
                    <td key={cell.id} className="border p-2">
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
          {/* Pagination */}
          <div className="flex items-center gap-2 mt-4">
            <button
              className="px-2 py-1 border rounded disabled:opacity-50"
              onClick={() => table.previousPage()}
              disabled={!table.getCanPreviousPage()}
            >
              ← 前へ
            </button>
            <button
              className="px-2 py-1 border rounded disabled:opacity-50"
              onClick={() => table.nextPage()}
              disabled={!table.getCanNextPage()}
            >
              次へ →
            </button>
            <span className="ml-2">
              Page{" "}
              <strong>
                {table.getState().pagination.pageIndex + 1} / {table.getPageCount()}
              </strong>
            </span>
            <select
              className="ml-2 border p-1"
              value={table.getState().pagination.pageSize}
              onChange={(e) => table.setPageSize(Number(e.target.value))}
            >
              {[5, 10, 20].map((pageSize) => (
                <option key={pageSize} value={pageSize}>
                  {pageSize} rows
                </option>
              ))}
            </select>
          </div>          


        </div>
      </div>

      <ItemDialog
        isOpen={dialogOpen}
        onClose={() => setDialogOpen(false)}
        onSave={handleSave}
        item={editingItem}
        mode={dialogMode}
      />
    </div>
  );
}

export default App;