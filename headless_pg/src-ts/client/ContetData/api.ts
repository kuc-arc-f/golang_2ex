import { Item, NewItem } from '../types/Item';
import dataUtil from './dataUtil';

const API_BASE = '/api/data'; 

export const itemsApi = {
  getAll: async (id: string): Promise<Item[]> => {
    const response = await fetch(`/api/admin/data_list?content=${id}`);
    if (!response.ok) {
      throw new Error('Failed to fetch items');
    }
    let json = await response.json();
    const items = dataUtil.getDataItems(json);
    //console.log(items);
    return items;
  },

  getById: async (id: number): Promise<Item> => {
    const response = await fetch(`${API_BASE}/${id}`);
    if (!response.ok) {
      throw new Error('Failed to fetch item');
    }
    return response.json();
  },

  create: async (item: NewItem): Promise<Item> => {
    const response = await fetch(API_BASE + "/create", {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(item),
    });
    if (!response.ok) {
      throw new Error('Failed to create item');
    }
    return response.json();
  },

  update: async (id: number, item: Partial<NewItem>): Promise<Item> => {
    item.id = id;
    const response = await fetch(`${API_BASE}/update`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(item),
    });
    if (!response.ok) {
      throw new Error('Failed to update item');
    }
    return response.json();
  },

  delete: async (id: number): Promise<void> => {
    const item = { id: id }
    const response = await fetch(`${API_BASE}/delete`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(item),
    });
    if (!response.ok) {
      throw new Error('Failed to delete item');
    }
  },
};