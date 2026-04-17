"use client";

import React, { useEffect, useState, useCallback } from "react";
import PageBreadcrumb from "@/components/common/PageBreadCrumb";
import ComponentCard from "@/components/common/ComponentCard";
import HtpasswdTable from "@/components/tables/HtpasswdTable";
import Button from "@/components/ui/button/Button";
import { Modal } from "@/components/ui/modal";
import Badge from "@/components/ui/badge/Badge";
import {
  listNamespaces,
  listHtpasswdSecrets,
  getHtpasswdSecret,
  createHtpasswdSecret,
  addHtpasswdUser,
  removeHtpasswdUser,
  deleteHtpasswdSecret,
} from "@/lib/api";
import type {
  HtpasswdSecretSummary,
  HtpasswdSecretDetail,
} from "@/lib/api";

export default function HtpasswdPage() {
  const [namespaces, setNamespaces] = useState<string[]>([]);
  const [selectedNamespace, setSelectedNamespace] = useState<string>("");

  const [secrets, setSecrets] = useState<HtpasswdSecretSummary[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [selectedName, setSelectedName] = useState<string | undefined>();
  const [detail, setDetail] = useState<HtpasswdSecretDetail | null>(null);
  const [detailLoading, setDetailLoading] = useState(false);

  const [createOpen, setCreateOpen] = useState(false);
  const [createName, setCreateName] = useState("");
  const [createUsername, setCreateUsername] = useState("");
  const [createPassword, setCreatePassword] = useState("");
  const [createSubmitting, setCreateSubmitting] = useState(false);

  const [addUserOpen, setAddUserOpen] = useState(false);
  const [addUsername, setAddUsername] = useState("");
  const [addPassword, setAddPassword] = useState("");
  const [addSubmitting, setAddSubmitting] = useState(false);

  const [deleteTarget, setDeleteTarget] = useState<{ type: "secret" | "user"; name: string; username?: string } | null>(null);
  const [deleteSubmitting, setDeleteSubmitting] = useState(false);

  const [actionError, setActionError] = useState<string | null>(null);

  const fetchSecrets = useCallback(async () => {
    if (!selectedNamespace) {
      setSecrets([]);
      setLoading(false);
      return;
    }
    try {
      setLoading(true);
      setError(null);
      const res = await listHtpasswdSecrets(selectedNamespace);
      setSecrets(res.data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load htpasswd secrets");
    } finally {
      setLoading(false);
    }
  }, [selectedNamespace]);

  useEffect(() => {
    listNamespaces()
      .then((res) => {
        setNamespaces(res.data);
        if (res.data.length > 0 && !selectedNamespace) {
          setSelectedNamespace(res.data[0]);
        }
      })
      .catch(() => {});
  }, []);

  useEffect(() => {
    fetchSecrets();
    setSelectedName(undefined);
    setDetail(null);
  }, [fetchSecrets]);

  const fetchDetail = useCallback(async (name: string) => {
    if (!selectedNamespace) return;
    try {
      setDetailLoading(true);
      setActionError(null);
      const res = await getHtpasswdSecret(selectedNamespace, name);
      setDetail(res.data);
    } catch (err) {
      setActionError(err instanceof Error ? err.message : "Failed to load secret detail");
    } finally {
      setDetailLoading(false);
    }
  }, [selectedNamespace]);

  useEffect(() => {
    if (selectedName) {
      fetchDetail(selectedName);
    } else {
      setDetail(null);
    }
  }, [selectedName, fetchDetail]);

  const handleSelectSecret = (name: string) => {
    setSelectedName(name);
  };

  const handleCreate = async () => {
    if (!selectedNamespace || !createName || !createUsername || !createPassword) return;
    try {
      setCreateSubmitting(true);
      setActionError(null);
      await createHtpasswdSecret(selectedNamespace, {
        name: createName,
        users: [{ username: createUsername, password: createPassword }],
      });
      setCreateOpen(false);
      setCreateName("");
      setCreateUsername("");
      setCreatePassword("");
      fetchSecrets();
    } catch (err) {
      setActionError(err instanceof Error ? err.message : "Failed to create secret");
    } finally {
      setCreateSubmitting(false);
    }
  };

  const handleAddUser = async () => {
    if (!selectedNamespace || !selectedName || !addUsername || !addPassword) return;
    try {
      setAddSubmitting(true);
      setActionError(null);
      await addHtpasswdUser(selectedNamespace, selectedName, {
        username: addUsername,
        password: addPassword,
      });
      setAddUserOpen(false);
      setAddUsername("");
      setAddPassword("");
      fetchDetail(selectedName);
      fetchSecrets();
    } catch (err) {
      setActionError(err instanceof Error ? err.message : "Failed to add user");
    } finally {
      setAddSubmitting(false);
    }
  };

  const handleDelete = async () => {
    if (!selectedNamespace || !deleteTarget) return;
    try {
      setDeleteSubmitting(true);
      setActionError(null);
      if (deleteTarget.type === "secret") {
        await deleteHtpasswdSecret(selectedNamespace, deleteTarget.name);
        setDeleteTarget(null);
        setSelectedName(undefined);
        setDetail(null);
        fetchSecrets();
      } else if (deleteTarget.username) {
        await removeHtpasswdUser(selectedNamespace, deleteTarget.name, deleteTarget.username);
        setDeleteTarget(null);
        fetchDetail(deleteTarget.name);
        fetchSecrets();
      }
    } catch (err) {
      setActionError(err instanceof Error ? err.message : "Failed to delete");
    } finally {
      setDeleteSubmitting(false);
    }
  };

  return (
    <div>
      <PageBreadcrumb pageTitle="Htpasswd Management" />
      <div className="space-y-6">
        <ComponentCard title="Htpasswd Secrets" desc="Manage .htpasswd credentials for BasicAuth">
          <div className="mb-4 flex items-center justify-between gap-3">
            <div className="flex items-center gap-3">
              <label
                htmlFor="namespace-filter"
                className="text-sm font-medium text-gray-700 dark:text-gray-300"
              >
                Namespace:
              </label>
              <select
                id="namespace-filter"
                value={selectedNamespace}
                onChange={(e) => setSelectedNamespace(e.target.value)}
                className="rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm text-gray-700 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300 focus:border-brand-500 focus:outline-none focus:ring-1 focus:ring-brand-500"
              >
                <option value="">Select Namespace</option>
                {namespaces.map((ns) => (
                  <option key={ns} value={ns}>
                    {ns}
                  </option>
                ))}
              </select>
            </div>
            {selectedNamespace && (
              <Button size="sm" onClick={() => setCreateOpen(true)}>
                Create Secret
              </Button>
            )}
          </div>

          {error && (
            <div className="mb-4 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700 dark:border-red-800 dark:bg-red-900/20 dark:text-red-400">
              {error}
            </div>
          )}

          {actionError && (
            <div className="mb-4 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700 dark:border-red-800 dark:bg-red-900/20 dark:text-red-400">
              {actionError}
            </div>
          )}

          <HtpasswdTable
            secrets={secrets}
            loading={loading}
            onSelect={handleSelectSecret}
            selectedName={selectedName}
          />
        </ComponentCard>

        {selectedName && (
          <ComponentCard title={`Secret: ${selectedName}`} desc="Users and linked SecurityPolicies">
            <div className="flex items-center gap-3 mb-4">
              <Button size="sm" onClick={() => setAddUserOpen(true)}>
                Add User
              </Button>
              <Button size="sm" variant="outline" onClick={() => setDeleteTarget({ type: "secret", name: selectedName })}>
                Delete Secret
              </Button>
            </div>

            {detailLoading ? (
              <p className="text-sm text-gray-500 dark:text-gray-400">Loading...</p>
            ) : detail ? (
              <div className="space-y-4">
                <div>
                  <h4 className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Users ({detail.userCount})
                  </h4>
                  {(detail.users ?? []).length === 0 ? (
                    <p className="text-sm text-gray-400">No users</p>
                  ) : (
                    <div className="flex flex-wrap gap-2">
                      {(detail.users ?? []).map((user) => (
                        <span
                          key={user.username}
                          className="inline-flex items-center gap-1 rounded-md bg-gray-100 px-3 py-1 text-sm text-gray-700 dark:bg-gray-800 dark:text-gray-300"
                        >
                          {user.username}
                          <button
                            onClick={() => setDeleteTarget({ type: "user", name: selectedName, username: user.username })}
                            className="ml-1 text-gray-400 hover:text-red-500 dark:hover:text-red-400"
                            title={`Remove ${user.username}`}
                          >
                            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                              <path d="M18 6L6 18M6 6l12 12" />
                            </svg>
                          </button>
                        </span>
                      ))}
                    </div>
                  )}
                </div>

                <div>
                  <h4 className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Linked SecurityPolicies
                  </h4>
                  {(detail.linkedSecurityPolicies ?? []).length === 0 ? (
                    <p className="text-sm text-gray-400">No linked policies</p>
                  ) : (
                    <div className="space-y-2">
                      {(detail.linkedSecurityPolicies ?? []).map((policy) => (
                        <div
                          key={policy.name}
                          className="flex items-center gap-2 rounded-lg border border-gray-200 px-3 py-2 dark:border-gray-700"
                        >
                          <Badge variant="light" color="primary" size="sm">
                            {policy.name}
                          </Badge>
                          <span className="text-xs text-gray-500 dark:text-gray-400">
                            {policy.targetRef.kind}: {policy.targetRef.name}
                          </span>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            ) : null}
          </ComponentCard>
        )}
      </div>

      {/* Create Secret Modal */}
      <Modal isOpen={createOpen} onClose={() => setCreateOpen(false)} className="max-w-[480px] p-6">
        <h3 className="mb-4 text-lg font-semibold text-gray-800 dark:text-white">Create Htpasswd Secret</h3>
        <div className="space-y-3">
          <div>
            <label className="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">Secret Name</label>
            <input
              type="text"
              value={createName}
              onChange={(e) => setCreateName(e.target.value)}
              className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300"
              placeholder="my-app-htpasswd"
            />
          </div>
          <div>
            <label className="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">Initial Username</label>
            <input
              type="text"
              value={createUsername}
              onChange={(e) => setCreateUsername(e.target.value)}
              className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300"
              placeholder="admin"
            />
          </div>
          <div>
            <label className="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">Password</label>
            <input
              type="password"
              value={createPassword}
              onChange={(e) => setCreatePassword(e.target.value)}
              className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300"
              placeholder="Min 8 characters"
            />
          </div>
          <div className="flex justify-end gap-2 pt-2">
            <Button size="sm" variant="outline" onClick={() => setCreateOpen(false)}>Cancel</Button>
            <Button size="sm" onClick={handleCreate} disabled={createSubmitting || !createName || !createUsername || !createPassword}>
              {createSubmitting ? "Creating..." : "Create"}
            </Button>
          </div>
        </div>
      </Modal>

      {/* Add User Modal */}
      <Modal isOpen={addUserOpen} onClose={() => setAddUserOpen(false)} className="max-w-[480px] p-6">
        <h3 className="mb-4 text-lg font-semibold text-gray-800 dark:text-white">Add User</h3>
        <div className="space-y-3">
          <div>
            <label className="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">Username</label>
            <input
              type="text"
              value={addUsername}
              onChange={(e) => setAddUsername(e.target.value)}
              className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300"
              placeholder="username"
            />
          </div>
          <div>
            <label className="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">Password</label>
            <input
              type="password"
              value={addPassword}
              onChange={(e) => setAddPassword(e.target.value)}
              className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300"
              placeholder="Min 8 characters"
            />
          </div>
          <div className="flex justify-end gap-2 pt-2">
            <Button size="sm" variant="outline" onClick={() => setAddUserOpen(false)}>Cancel</Button>
            <Button size="sm" onClick={handleAddUser} disabled={addSubmitting || !addUsername || !addPassword}>
              {addSubmitting ? "Adding..." : "Add User"}
            </Button>
          </div>
        </div>
      </Modal>

      {/* Delete Confirmation Modal */}
      <Modal isOpen={deleteTarget !== null} onClose={() => setDeleteTarget(null)} className="max-w-[400px] p-6">
        <h3 className="mb-4 text-lg font-semibold text-gray-800 dark:text-white">Confirm Delete</h3>
        <p className="mb-4 text-sm text-gray-600 dark:text-gray-400">
          {deleteTarget?.type === "secret"
            ? `Delete secret "${deleteTarget.name}" and all its users?`
            : `Remove user "${deleteTarget?.username}" from "${deleteTarget?.name}"?`}
        </p>
        <div className="flex justify-end gap-2">
          <Button size="sm" variant="outline" onClick={() => setDeleteTarget(null)}>Cancel</Button>
          <Button size="sm" onClick={handleDelete} disabled={deleteSubmitting}>
            {deleteSubmitting ? "Deleting..." : "Delete"}
          </Button>
        </div>
      </Modal>
    </div>
  );
}
