"use client";

import React, { useEffect, useState, useCallback } from "react";
import PageBreadcrumb from "@/components/common/PageBreadCrumb";
import ComponentCard from "@/components/common/ComponentCard";
import ServiceTable from "@/components/tables/ServiceTable";
import { listServices, listNamespaces } from "@/lib/api";
import type { ServiceOverview } from "@/lib/api";

export default function ServicesPage() {
  const [services, setServices] = useState<ServiceOverview[]>([]);
  const [namespaces, setNamespaces] = useState<string[]>([]);
  const [selectedNamespace, setSelectedNamespace] = useState<string>("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchServices = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const ns = selectedNamespace || undefined;
      const res = await listServices(ns);
      setServices(res.data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load services");
    } finally {
      setLoading(false);
    }
  }, [selectedNamespace]);

  useEffect(() => {
    listNamespaces()
      .then((res) => setNamespaces(res.data))
      .catch(() => {
        // namespaces are optional filter, ignore error
      });
  }, []);

  useEffect(() => {
    fetchServices();
  }, [fetchServices]);

  return (
    <div>
      <PageBreadcrumb pageTitle="Service Overview" />
      <div className="space-y-6">
        <ComponentCard
          title="Kubernetes Services"
          desc="Services with HTTPRoute and SecurityPolicy status"
        >
          <div className="mb-4 flex items-center gap-3">
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
              <option value="">All Namespaces</option>
              {namespaces.map((ns) => (
                <option key={ns} value={ns}>
                  {ns}
                </option>
              ))}
            </select>
          </div>

          {error && (
            <div className="mb-4 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700 dark:border-red-800 dark:bg-red-900/20 dark:text-red-400">
              {error}
            </div>
          )}

          <ServiceTable services={services} loading={loading} />
        </ComponentCard>
      </div>
    </div>
  );
}
