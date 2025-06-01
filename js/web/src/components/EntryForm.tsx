import React, { useState } from 'react';
import { create } from '@bufbuild/protobuf';
import { entryClient } from '../api/client';
import { CreateEntryRequestSchema, GrowthStage } from '../genproto/protobuf/entry/entry_pb';

interface EntryFormProps {
  onEntryCreated?: (entryId: string) => void;
}

const EntryForm: React.FC<EntryFormProps> = ({ onEntryCreated }) => {
  const [formData, setFormData] = useState({
    title: '',
    content: '',
    growthStage: GrowthStage.SEED,
    userId: '' // Add this in case it's required
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string>('');
  const [success, setSuccess] = useState<string>('');

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: name === 'growthStage' ? parseInt(value) : value
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setError('');
    setSuccess('');

    try {
      // Create the request with form data
      const request = create(CreateEntryRequestSchema, {
        title: formData.title.trim(),
        content: formData.content.trim(),
        growthStage: formData.growthStage,
        ...(formData.userId.trim() && { userId: formData.userId.trim() }) // Only add if not empty
      });

      console.log('Sending request:', request);
      
      const response = await entryClient.createEntry(request);
      
      if (response.entry?.id) {
        setSuccess(`✅ Entry created successfully! ID: ${response.entry.id}`);
        setFormData({
          title: '',
          content: '',
          growthStage: GrowthStage.SEED,
          userId: ''
        });
        onEntryCreated?.(response.entry.id);
      }
      
    } catch (error: any) {
      console.error('Failed to create entry:', error);
      setError(`❌ Failed to create entry: ${error.message || 'Unknown error'}`);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="entry-form">
      <h2>Create New Entry</h2>
      
      {error && (
        <div className="alert alert-error">
          {error}
        </div>
      )}
      
      {success && (
        <div className="alert alert-success">
          {success}
        </div>
      )}

      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="userId">User ID (optional):</label>
          <input
            type="text"
            id="userId"
            name="userId"
            value={formData.userId}
            onChange={handleInputChange}
            placeholder="Leave empty if not required"
          />
        </div>

        <div className="form-group">
          <label htmlFor="title">Title *:</label>
          <input
            type="text"
            id="title"
            name="title"
            value={formData.title}
            onChange={handleInputChange}
            required
            placeholder="Enter entry title"
          />
        </div>

        <div className="form-group">
          <label htmlFor="content">Content *:</label>
          <textarea
            id="content"
            name="content"
            value={formData.content}
            onChange={handleInputChange}
            required
            placeholder="Enter entry content"
            rows={4}
          />
        </div>

        <div className="form-group">
          <label htmlFor="growthStage">Growth Stage:</label>
          <select
            id="growthStage"
            name="growthStage"
            value={formData.growthStage}
            onChange={handleInputChange}
          >
            <option value={GrowthStage.SEED}>Seed</option>
            <option value={GrowthStage.SPROUT}>Sprout</option>
            <option value={GrowthStage.BLOOM}>Bloom</option>
          </select>
        </div>

        <button 
          type="submit" 
          disabled={isSubmitting || !formData.title.trim() || !formData.content.trim()}
          className="btn-primary"
        >
          {isSubmitting ? 'Creating...' : 'Create Entry'}
        </button>
      </form>
    </div>
  );
};

export default EntryForm;